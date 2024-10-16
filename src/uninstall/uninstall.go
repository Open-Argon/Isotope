package uninstall

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/help"
	"github.com/Open-Argon/Isotope/src/package/zip"
)

const usage = `uninstall [package] [options]`

var o = help.Options{
	{"install globally or locally into your CWD", "--global, -g"},
	{"show help", "--help, -h"},
}

func deleteDirectoryIfExists(dirPath string) error {
	// Check if the directory exists
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Directory does not exist, return without error
			return nil
		}
		// Error occurred while checking directory existence
		return err
	}

	// Delete the directory
	err = os.RemoveAll(dirPath)
	if err != nil {
		return err
	}

	return nil
}

func Uninstall() {
	var args = args.Args[1:]
	var name = ""
	var global = false
	var deleted = false
	var path = ""

	if len(args) == 0 {
		help.Help(usage, o)
		return
	}

	if len(args) > 0 {
		if args[0] == "--help" || args[0] == "-h" {
			help.Help(usage, o)
			return
		}
		name = args[0]
	}
	if len(args) > 1 {
		if args[1] == "--global" || args[1] == "-g" {
			global = true
		}
	}
	var err error
	if global {
		path, err = os.Executable()
		path = filepath.Dir(path)
	} else {
		path, err = os.Getwd()
	}
	if err != nil {
		log.Fatal(err)
	}
	if global {
		var delpath = filepath.Join(path, "argon_modules", name)
		if _, err := os.Stat(delpath); os.IsNotExist(err) {
			log.Fatal("Package not found")
		}
		deleteDirectoryIfExists(delpath)
	} else {
		lockPath := filepath.Join(path, "iso-lock.json")
		lockFile, err := os.Open(lockPath)
		if err != nil {
			log.Fatal(err)
		}
		var dependencies = make([]zip.Dependency, 0)
		err = json.NewDecoder(lockFile).Decode(&dependencies)
		if err != nil {
			log.Fatal(err)
		}
		lockFile.Close()
		var newDependencies = make([]zip.Dependency, 0)
		for index, dependency := range dependencies {
			if dependency.Name == name {
				deleteDirectoryIfExists(filepath.Join(path, "argon_modules", dependency.Name))
				deleted = true
			} else {
				newDependencies = append(newDependencies, dependencies[index])
			}
		}
		if !deleted {
			log.Fatal("Package not found")
		}
		lockFile, err = os.Create(lockPath)
		if err != nil {
			log.Fatal(err)
		}
		err = json.NewEncoder(lockFile).Encode(newDependencies)
		if err != nil {
			log.Fatal(err)
		}
		lockFile.Close()
	}
	fmt.Println("Uninstalled package:", name)
}
