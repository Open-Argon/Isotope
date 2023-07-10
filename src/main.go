package main

import (
	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/help"
	Init "github.com/Open-Argon/Isotope/src/init"
	"github.com/Open-Argon/Isotope/src/install"
	pack "github.com/Open-Argon/Isotope/src/package"
	"github.com/Open-Argon/Isotope/src/version"
)

var usage = "isotope <command> [options]"
var options = help.Options{
	{"install a package", "install [options]"},
	{"package a project", "package [options]"},
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
	case "package":
		pack.Pack()
	case "help":
		help.Help(usage, options)
	case "version":
		version.PrintVersion()
	case "init":
		Init.Init()
	}
}
