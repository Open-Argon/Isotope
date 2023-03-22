package install

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/help"
)

var usage = `install [options] [package]`
var o = help.Options{
	{"specify a specific remote host", "--remote [host]"},
	{"show help", "--help, -h"},
}

func Install() {
	var args = args.Args[1:]
	if len(args) == 0 {
		help.Help(usage, o)
		return
	}
	name := args[0]
	remote := "https://pkg.argon.wbell.dev/"
	version := "latest"
	switch args[0] {
	case "--help", "-h":
		help.Help(usage, o)
		return
	case "--remote":
		if len(args) < 3 {
			help.Help(usage, o)
			return
		}
		remote = args[1]
		name = args[2]
	}
	fmt.Println("Searching for", name+"@"+version)
	fmt.Println("Remote:", remote)
	time.Sleep(1 * time.Second)
	fmt.Print("\nDownloading")
	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(rand.Float64()*250) * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Print("\n\nInstalling")
	for i := 0; i < 100; i++ {
		time.Sleep(time.Duration(rand.Float64()*100) * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Print("\n\nInstalled!\n")
}
