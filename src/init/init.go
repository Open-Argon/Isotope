package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	validpackagename "github.com/Open-Argon/Isotope/src/validPackageName"
)

func Init() {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the current directory name
	dirName := filepath.Base(dir)

	packageName := validpackagename.ValidPackageName(dirName)

	// Create the package.json file
	packageFile, err := os.Create("iso-package.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer packageFile.Close()

	// Write the package.json file
	packageFile.WriteString(`{
	"name": "` + packageName + `",
	"version": "1.0.0"
}`)

	// Create the package-lock.json file
	packageLockFile, err := os.Create("iso-package-lock.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer packageLockFile.Close()

	// Write the package-lock.json file
	packageLockFile.WriteString(`[]`)
	fmt.Println("Initialized empty Isotope package in", strings.ReplaceAll(dir, "\\", "/"))
}
