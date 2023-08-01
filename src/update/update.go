package update

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Open-Argon/Isotope/src/config"
	"github.com/Open-Argon/Isotope/src/install"
	"github.com/Open-Argon/Isotope/src/package/zip"
)

func Update() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	lockfilePath := filepath.Join(cwd, "iso-lock.json")
	lockfile, err := os.Open(lockfilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer lockfile.Close()
	dependencies := make([]zip.Dependency, 0)
	err = json.NewDecoder(lockfile).Decode(&dependencies)
	if err != nil {
		log.Fatal(err)
	}
	outputDependencies := make([]zip.Dependency, len(dependencies))
	for i, dependency := range dependencies {
		var remote = dependency.Remote
		if remote == "" {
			remote = config.DefaultRemote
		}
		pkg := install.InstallPackage(remote, "", dependency.Name, "latest", cwd, make([]zip.Dependency, 0))
		if pkg.Version != dependency.Version {
			fmt.Println("Updated", dependency.Name+"@"+dependency.Version, ">>", pkg.Name+"@"+pkg.Version)
		} else {
			fmt.Println("Up to date", pkg.Name+"@"+pkg.Version)
		}
		outputDependencies[i] = pkg
	}
	lockfile, err = os.Create(lockfilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer lockfile.Close()
	err = json.NewEncoder(lockfile).Encode(outputDependencies)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
}
