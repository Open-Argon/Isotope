package install

import (
	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/help"
)

var usage = `Usage: isotope install [options] [package]`
var o = help.Options{
	{"show help", "--help, -h, help"},
}

func Install() {
	var args = args.Args[1:]
	if len(args) == 0 {
		return
	}
	switch args[0] {
	case "--help", "-h", "help":
		help.Help(usage, o)
	}
}
