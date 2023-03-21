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
	"github.com/Open-Argon/Isotope/src/start"
)

var usage = `zip [options] [path]`
var o = help.Options{
	{"Show this help message", "--help, -h"},
}

func Zip() {
	args := args.Args[2:]
	if len(args) == 0 {
		help.Help(usage, o)
		return
	}
	p := args[0]
	switch args[0] {
	case "--help", "-h":
		help.Help(usage, o)
		return
	}

	fmt.Println("zipping package...")
	packageFile, buf := ReadPackageAndDependencies(p)
	hash := hash.Sha256Hex(packageFile["name"].(string))
	zipPath := path.Join(start.ZippedPath, hash+".zip")
	fmt.Println("done")
	fmt.Println()
	fmt.Println("Package:", packageFile["name"])
	fmt.Println("Version:", packageFile["version"])
	fmt.Println(len(packageFile["dependencies"].([]any)), "dependencies")
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
