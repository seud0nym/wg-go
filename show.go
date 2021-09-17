package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func show(opts *cmdOptions) {
	if opts.Interface == "--help" || (opts.Interface == "interfaces" && opts.Option != "") || !(opts.Option == "" || opts.Option == "public-key" || opts.Option == "private-key" || opts.Option == "listen-port" || opts.Option == "fwmark" || opts.Option == "peers" || opts.Option == "preshared-keys" || opts.Option == "endpoints" || opts.Option == "allowed-ips" || opts.Option == "latest-handshakes" || opts.Option == "transfer" || opts.Option == "persistent-keepalive" || opts.Option == "dump") {
		showSubCommandUsage("show { <interface> | all | interfaces } [public-key | private-key | listen-port | fwmark | peers | preshared-keys | endpoints | allowed-ips | latest-handshakes | transfer | persistent-keepalive | dump]", opts)
	}

	client, err := wgctrl.New()
	checkError(err)
	switch opts.Interface {
	case "interfaces":
		devices, err := client.Devices()
		checkError(err)
		for i := 0; i < len(devices); i++ {
			fmt.Println(devices[i].Name)
		}
		break
	case "all":
		devices, err := client.Devices()
		checkError(err)
		for _, dev := range devices {
			showDevice(*dev, opts)
		}
		break
	default:
		dev, err := client.Device(opts.Interface)
		checkError(err)
		showDevice(*dev, opts)
	}
	client.Close()
}

func showDevice(dev wgtypes.Device, opts *cmdOptions) {
	if opts.Option == "" {
		showKeys := opts.ShowKeys
		fmt.Printf("Interface: %s (%s)\n", dev.Name, dev.Type.String())
		fmt.Printf("  public key: %s\n", dev.PublicKey.String())
		fmt.Printf("  private key: %s\n", formatKey(dev.PrivateKey, showKeys))
		fmt.Printf("  listening port: %d\n", dev.ListenPort)
		fmt.Println()
		for _, peer := range dev.Peers {
			showPeers(peer, showKeys)
		}
	} else {
		deviceName := ""
		if opts.Interface == "all" {
			deviceName = dev.Name + "\t"
		}
		switch opts.Option {
		case "public-key":
			fmt.Printf("%s%s\n", deviceName, dev.PublicKey.String())
			break
		case "private-key":
			fmt.Printf("%s%s\n", deviceName, dev.PrivateKey.String())
			break
		case "listen-port":
			fmt.Printf("%s%d\n", deviceName, dev.ListenPort)
			break
		case "fwmark":
			fmt.Printf("%s%d\n", deviceName, dev.FirewallMark)
			break
		case "peers":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\n", deviceName, peer.PublicKey.String())
			}
			break
		case "preshared-keys":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), formatPSK(peer.PresharedKey, "(none)"))
			}
			break
		case "endpoints":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), peer.Endpoint.String())
			}
			break
		case "allowed-ips":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), joinIPs(peer.AllowedIPs))
			}
			break
		case "latest-handshakes":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%d\n", deviceName, peer.PublicKey.String(), peer.LastHandshakeTime.Unix())
			}
			break
		case "transfer":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%d\t%d\n", deviceName, peer.PublicKey.String(), peer.ReceiveBytes, peer.TransmitBytes)
			}
			break
		case "persistent-keepalive":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), zeroToOff(strconv.FormatFloat(peer.PersistentKeepaliveInterval.Seconds(), 'g', 0, 64)))
			}
			break
		case "dump":
			fmt.Printf("%s%s\t%s\t%d\t%s\n", deviceName, dev.PrivateKey.String(), dev.PublicKey.String(), dev.ListenPort, zeroToOff(strconv.FormatInt(int64(dev.FirewallMark), 10)))
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\t%s\t%s\t%d\t%d\t%d\t%s\n",
					deviceName,
					peer.PublicKey.String(),
					formatPSK(peer.PresharedKey, "(none)"),
					peer.Endpoint.String(),
					joinIPs(peer.AllowedIPs),
					peer.LastHandshakeTime.Unix(),
					peer.ReceiveBytes,
					peer.TransmitBytes,
					zeroToOff(strconv.FormatFloat(peer.PersistentKeepaliveInterval.Seconds(), 'g', 0, 64)))
			}
			break
		}
	}
}

func showPeers(peer wgtypes.Peer, showKeys bool) {
	const tmpl = `peer: {{ .PublicKey }}
  endpoint = {{ .Endpoint }}
  allowed ips = {{ .AllowedIPs }}
  {{- if .PresharedKey}}
  preshared key = {{ .PresharedKey }}
  {{- end}}
  keep alive interval = {{ .KeepAliveInterval }}s
  last handshake time = {{ .LastHandshakeTime }}
  transfer: {{ .ReceiveBytes }} bytes received, {{ .TransmitBytes }} bytes sent
  protocol version = {{ .ProtocolVersion }} 

`
	type tmplContent struct {
		PublicKey         string
		PresharedKey      string
		Endpoint          string
		KeepAliveInterval float64
		LastHandshakeTime string
		ReceiveBytes      int64
		TransmitBytes     int64
		AllowedIPs        string
		ProtocolVersion   int
	}

	t := template.Must(template.New("peer_tmpl").Parse(tmpl))
	c := tmplContent{
		PublicKey:         peer.PublicKey.String(),
		PresharedKey:      formatPSK(peer.PresharedKey, ""),
		Endpoint:          peer.Endpoint.String(),
		KeepAliveInterval: peer.PersistentKeepaliveInterval.Seconds(),
		LastHandshakeTime: peer.LastHandshakeTime.Format(time.RFC3339),
		ReceiveBytes:      peer.ReceiveBytes,
		TransmitBytes:     peer.TransmitBytes,
		AllowedIPs:        joinIPs(peer.AllowedIPs),
		ProtocolVersion:   peer.ProtocolVersion,
	}

	err := t.Execute(os.Stdout, c)
	checkError(err)
}

func formatKey(key wgtypes.Key, showKeys bool) (string) {
	k := "(hidden)"
	if showKeys {
		k = key.String()
	}
	return k
}

func formatPSK(key wgtypes.Key, none string) (string) {
	psk := key.String()
	if psk == "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" {
		return none
	}
	return psk
}

func joinIPs(ips []net.IPNet) (string) {
	ipStrings := make([]string, 0, len(ips))
	for _, v := range ips {
		ipStrings = append(ipStrings, v.String())
	}
	return strings.Join(ipStrings, ", ")
}

func zeroToOff(value string) (string) {
	if value == "0" {
		return "off"
	}
	return value
}
