package main

import (
	"fmt"
	"log"
	"os"
)

const (
	ENV_WG_COMMAND = "WG_COMMAND"
)

var (
	appVersion = "dev"
)

func main() {
	opts := parseArgs()

	switch opts.SubCommand {
	case "show":
		show(opts)
		break
	case "showconf":
		showConfig(opts)
		break
	case "setconf":
		setConfig(opts)
		break
	case "genkey":
		genKey(opts)
		break
	case "genpsk":
		genPSK(opts)
		break
	case "pubkey":
		pubKey(opts)
		break
	case "version":
		showVersion(opts)
		break
	default:
		fmt.Printf("Invalid subcommand: '%s'\n", opts.Command)
		showCommandUsage(1, opts)
	}
}

type cmdOptions struct {
	Command    string
	SubCommand string
	Interface  string
	Option     string
}

func parseArgs() *cmdOptions {
	args := len(os.Args[1:])
	base := 0
	opts := cmdOptions{}

	opts.Command = os.Getenv(ENV_WG_COMMAND)
	if opts.Command == "" {
		opts.Command = "wg-go"
	}

	if args == 0 {
		opts.SubCommand = "show"
		opts.Interface = "all"
	} else if args == 1 && os.Args[base+1] == "--help" {
		showCommandUsage(0, &opts)
	} else if args > 3 {
		showCommandUsage(1, &opts)
	} else {
		opts.SubCommand = os.Args[base+1]
		if args >= 2 {
			opts.Interface = os.Args[base+2]
			if args == 3 {
				opts.Option = os.Args[base+3]
			}
		} else if opts.SubCommand == "show" {
			opts.Interface = "all"
		}
	}

	return &opts
}

func showCommandUsage(code int, opts *cmdOptions) {
	subcommands := `Available subcommands:
  show:     Shows the current configuration and device information
  showconf: Shows the current configuration of a given WireGuard interface, for use with 'setconf'
  setconf:  Applies a configuration file to a WireGuard interface
  genkey:   Generates a new private key and writes it to stdout
  genpsk:   Generates a new preshared key and writes it to stdout
  pubkey:   Reads a private key from stdin and writes a public key to stdout
  version:  Shows the version`

	fmt.Printf("Usage: %s <cmd> [<args>]\n\n", opts.Command)
	fmt.Printf("%s\n\n", subcommands)
	fmt.Println("You may pass '--help' to any of these subcommands to view showCommandUsage.")
	os.Exit(code)
}

func showSubCommandUsage(parameters string, opts *cmdOptions) {
	fmt.Printf("Usage: %s %s\n", opts.Command, parameters)
	if opts.Interface == "--help" {
		os.Exit(0)
	} else {
		os.Exit(2)
	}
}

func showVersion(opts *cmdOptions) {
	if opts.Interface == "--help" || opts.Interface != "" || opts.Option != "" {
		showSubCommandUsage("version", opts)
	}
	fmt.Printf("wg-go v%s https://github.com/seud0nym/wg-go\n", appVersion)
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
		os.Exit(2)
	}
}
