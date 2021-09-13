package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func genKey(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface != "" || opts.Option != "" {
		showSubCommandUsage("genkey", opts)
	}

	key, err := wgtypes.GeneratePrivateKey()
	checkError(err)
	fmt.Println(key.String())
}

func genPSK(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface != "" || opts.Option != "" {
		showSubCommandUsage("genpsk", opts)
	}

	key, err := wgtypes.GenerateKey()
	checkError(err)
	fmt.Println(key.String())
}

func pubKey(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface != "" || opts.Option != "" {
		showSubCommandUsage("pubkey", opts)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	checkError(err)
	input = strings.TrimSpace(input)
	private, err := wgtypes.ParseKey(input)
	checkError(err)
	public := private.PublicKey()
	fmt.Println(public.String())
}
