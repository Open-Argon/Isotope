package pack

import (
	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/help"
	"github.com/Open-Argon/Isotope/src/package/zip"
)

var usage = `pack [options]`
var o = help.Options{
	{"pack an module into a package", "zip [path]"},
	{"push a module to a remote", "push [package name]"},
	{"list all zipped packages", "list"},
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
	case "zip":
		zip.Zip()
	case "push":
		push()
	}
}
