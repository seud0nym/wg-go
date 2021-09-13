package main

import (
	"fmt"
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
		fmt.Printf("Interface: %s (%s)\n", dev.Name, dev.Type.String())
		fmt.Printf("  public key: %s\n", dev.PublicKey.String())
		fmt.Println("  private key: (hidden)")
		fmt.Printf("  listening port: %d\n", dev.ListenPort)
		for _, peer := range dev.Peers {
			showPeers(peer)
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
				psk := peer.PresharedKey.String()
				if psk == "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" {
					psk = "(none)"
				}
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), psk)
			}
			break
		case "endpoints":
			for _, peer := range dev.Peers {
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), peer.Endpoint.String())
			}
			break
		case "allowed-ips":
			for _, peer := range dev.Peers {
				allowdIpStrings := make([]string, 0, len(peer.AllowedIPs))
				for _, v := range peer.AllowedIPs {
					allowdIpStrings = append(allowdIpStrings, v.String())
				}
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), strings.Join(allowdIpStrings, ", "))
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
				ka := strconv.FormatFloat(peer.PersistentKeepaliveInterval.Seconds(), 'g', 0, 64)
				if ka == "0" {
					ka = "off"
				}
				fmt.Printf("%s%s\t%s\n", deviceName, peer.PublicKey.String(), ka)
			}
			break
		case "dump":
			fmark := strconv.FormatInt(int64(dev.FirewallMark), 10)
			if fmark == "0" {
				fmark = "off"
			}
			fmt.Printf("%s%s\t%s\t%d\t%s\n", deviceName, dev.PrivateKey.String(), dev.PublicKey.String(), dev.ListenPort, fmark)
			for _, peer := range dev.Peers {
				psk := peer.PresharedKey.String()
				if psk == "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" {
					psk = "(none)"
				}
				allowdIpStrings := make([]string, 0, len(peer.AllowedIPs))
				for _, v := range peer.AllowedIPs {
					allowdIpStrings = append(allowdIpStrings, v.String())
				}
				ka := strconv.FormatFloat(peer.PersistentKeepaliveInterval.Seconds(), 'f', 0, 64)
				if ka == "0" {
					ka = "off"
				}
				fmt.Printf("%s%s\t%s\t%s\t%s\t%d\t%d\t%d\t%s\n", deviceName, peer.PublicKey.String(), psk, peer.Endpoint.String(), strings.Join(allowdIpStrings, ", "), peer.LastHandshakeTime.Unix(), peer.ReceiveBytes, peer.TransmitBytes, ka)
			}
			break
		}
	}
}

func showPeers(peer wgtypes.Peer) {
	const tmpl = `
peer: {{ .PublicKey }}
  endpoint = {{ .Endpoint }}
  allowed ips = {{ .AllowedIPs }}
  preshared key = {{ .PresharedKey }}
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
		PresharedKey:      "(hidden)",
		Endpoint:          peer.Endpoint.String(),
		KeepAliveInterval: peer.PersistentKeepaliveInterval.Seconds(),
		LastHandshakeTime: peer.LastHandshakeTime.Format(time.RFC3339),
		ReceiveBytes:      peer.ReceiveBytes,
		TransmitBytes:     peer.TransmitBytes,
		AllowedIPs:        "",
		ProtocolVersion:   peer.ProtocolVersion,
	}

	allowdIpStrings := make([]string, 0, len(peer.AllowedIPs))
	for _, v := range peer.AllowedIPs {
		allowdIpStrings = append(allowdIpStrings, v.String())
	}
	c.AllowedIPs = strings.Join(allowdIpStrings, ", ")
	err := t.Execute(os.Stdout, c)
	checkError(err)
}
