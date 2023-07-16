package pack

import (
	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/help"
	"github.com/Open-Argon/Isotope/src/package/zip"
)

var usage = `pack [options]`
var o = help.Options{
	{"pack an module into a package", "build [path]"},
	{"push a module to a remote", "push"},
	{"show help", "--help, -h"},
}

func Pack() {
	args := args.Args[1:]
	if len(args) == 0 {
		help.Help(usage, o)
		return
	}
	switch args[0] {
	case "--help", "-h":
		help.Help(usage, o)
		return
	case "build":
		zip.Zip()
		return
	case "push":
		push()
		return
	}
	help.Help(usage, o)
}
