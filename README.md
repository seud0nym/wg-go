# wg-go

A Golang implementation of the WireGuard [wg(8)](https://git.zx2c4.com/wireguard-tools/about/src/man/wg.8) utility.

This tool could be used to get and set the configuration of WireGuard tunnel interfaces.

It can be used in conjunction with [wireguard-go](https://git.zx2c4.com/wireguard-go/about/) for an almost complete userspace implementation of Wireguard on platforms which can be targeted by Go but do not have an implementation of Wireguard available.

`wg-go` can also control a kernel-based Wireguard configuration.

For more information on WireGuard, please see https://www.wireguard.com/.

## Supported Sub-commands

This implementation supports the following sub-commands as specified in [wg(8)](https://git.zx2c4.com/wireguard-tools/about/src/man/wg.8):
```
  show:     Shows the current configuration and device information
  showconf: Shows the current configuration of a given WireGuard interface, for use with 'setconf'
  setconf:  Applies a configuration file to a WireGuard interface
  genkey:   Generates a new private key and writes it to stdout
  genpsk:   Generates a new preshared key and writes it to stdout
  pubkey:   Reads a private key from stdin and writes a public key to stdout
```

## How does this work?

This tool uses [wgctrl-go](https://github.com/WireGuard/wgctrl-go/) to enable control of Wireguard devices on multiple platforms.

## Original Code

This project was inspired by and based upon [QuantumGhost/wg-quick-go](https://github.com/QuantumGhost/wg-quick-go).
