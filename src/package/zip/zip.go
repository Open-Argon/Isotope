package zip

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/hash"
	"github.com/Open-Argon/Isotope/src/help"
)

var usage = `zip [path] [options]`
var o = help.Options{
	{"Show this help message", "--help, -h"},
}

func Zip() {
	args := args.Args[2:]
	p := ""
	if len(args) == 0 {
		var err error
		p, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		p = args[0]
		switch args[0] {
		case "--help", "-h":
			help.Help(usage, o)
			return
		}
	}

	fmt.Println("building package...")
	packageFile, buf := ReadPackageAndDependencies(p)
	hash := hash.Sha256Hex(packageFile.Name + "@" + packageFile.Version)
	zipPath := path.Join(p, "__isotope__", "builds", "armod-"+hash+".tar.gz")
	os.MkdirAll(path.Dir(zipPath), os.ModePerm)
	fmt.Println()
	fmt.Println("Package:", packageFile.Name)
	fmt.Println("Version:", packageFile.Version)
	fmt.Println(len(packageFile.Dependencies), "dependencies")
	fmt.Println("Size:", buf.Len(), "bytes")
	fmt.Println()
	fmt.Println("Saving to", strconv.Quote(zipPath))
	fmt.Println()

	err := os.WriteFile(zipPath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done!")
}
