package main

import (
	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/help"
	"github.com/Open-Argon/Isotope/src/install"
	"github.com/Open-Argon/Isotope/src/version"
)

var usage = "isotope <command> [options]"
var options = help.Options{
	{"install a package", "install [options]"},
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
	case "help":
		help.Help(usage, options)
	case "version":
		version.PrintVersion()
	}
}
