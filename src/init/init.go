package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	validpackagename "github.com/Open-Argon/Isotope/src/validPackageName"
)

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func Init(silent bool, dir string) {

	// Get the current directory name
	dirName := filepath.Base(dir)

	packageName := validpackagename.ValidPackageName(dirName)
	if packageName == "" {
		fmt.Println("Invalid package name:", dirName)
		os.Exit(1)
	}

	// Create the package.json file
	packageFilePath := filepath.Join(dir, "argon-package.json")
	if fileExists(packageFilePath) {
		if !silent {
			fmt.Println("argon-package.json already initialized in", strings.ReplaceAll(dir, "\\", "/"))
		}
	} else {
		packageFile, err := os.Create(packageFilePath)
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
	}

	// Create the package-lock.json file
	packageLockPath := filepath.Join(dir, "iso-lock.json")
	if fileExists(packageLockPath) {
		if !silent {
			fmt.Println("iso-lock.json already initialized in", strings.ReplaceAll(dir, "\\", "/"))
		}
	} else {
		packageLockFile, err := os.Create(packageLockPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer packageLockFile.Close()

		// Write the package-lock.json file
		packageLockFile.WriteString(`[]`)
		if !silent {
			fmt.Println("Initialized empty Isotope package in", strings.ReplaceAll(dir, "\\", "/"))
		}
	}
}
