package main

import (
	"fmt"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func showConfig(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface == "" {
		showSubCommandUsage("showconf <interface>", opts)
	}

	client, err := wgctrl.New()
	checkError(err)
	dev, err := client.Device(opts.Interface)
	checkError(err)
	fmt.Printf("[Interface]\n")
	fmt.Printf("ListenPort =  %d\n", dev.ListenPort)
	fmt.Printf("PrivateKey = %s\n", dev.PrivateKey.String())
	for _, peer := range dev.Peers {
		showConfigPeers(peer)
	}
}

func showConfigPeers(peer wgtypes.Peer) {
	allowdIpStrings := make([]string, 0, len(peer.AllowedIPs))
	for _, v := range peer.AllowedIPs {
		allowdIpStrings = append(allowdIpStrings, v.String())
	}
	psk := peer.PresharedKey.String()
	ka := peer.PersistentKeepaliveInterval.Seconds()

	fmt.Printf("\n[Peer]\n")
	fmt.Printf("PublicKey = %s\n", peer.PublicKey.String())
	if psk != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=" {
		fmt.Printf("PresharedKey = %s\n", peer.PresharedKey.String())
	}
	fmt.Printf("AllowedIPs = %s\n", strings.Join(allowdIpStrings, ", "))
	fmt.Printf("Endpoint = %s\n", peer.Endpoint.String())
	if ka > 0 {
		fmt.Printf("PersistentKeepalive = %g\n", ka)
	}
}
