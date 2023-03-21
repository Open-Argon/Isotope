package install

import (
	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/help"
)

var usage = `isotope install [options] [package]`
var o = help.Options{
	{"show help", "--help, -h"},
}

func Install() {
	var args = args.Args[1:]
	if len(args) == 0 {
		help.Help(usage, o)
		return
	}
	switch args[0] {
	case "--help", "-h":
		help.Help(usage, o)
	}
}
