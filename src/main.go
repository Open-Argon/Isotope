package main

import (
	"fmt"
	"os"

	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/auth"
	"github.com/Open-Argon/Isotope/src/help"
	Init "github.com/Open-Argon/Isotope/src/init"
	"github.com/Open-Argon/Isotope/src/install"
	pack "github.com/Open-Argon/Isotope/src/package"
	"github.com/Open-Argon/Isotope/src/uninstall"
	"github.com/Open-Argon/Isotope/src/update"
	"github.com/Open-Argon/Isotope/src/version"
)

var usage = "isotope <command> [options]"
var options = help.Options{
	{"install a package", "install [options]"},
	{"uninstall a package", "uninstall [options]"},
	{"package a project", "package [options]"},
	{"set auth token for pushing to a remote host", "auth [options]"},
	{"update all packages in the CWD", "update"},
	{"initialize a project", "init"},
	{"show help", "help"},
	{"show version", "version"},
}

func main() {
	if len(args.Args) == 0 {
		help.Help(usage, options)
		return
	}
	switch args.Args[0] {
	case "install":
		install.Install()
	case "uninstall":
		uninstall.Uninstall()
	case "package":
		pack.Pack()
	case "version":
		version.PrintVersion()
	case "init":
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		Init.Init(false, dir)
	case "auth":
		auth.Auth()
	case "update":
		update.Update()
	default:
		help.Help(usage, options)
	}
}
