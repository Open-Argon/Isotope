package install

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/config"
	"github.com/Open-Argon/Isotope/src/help"
	"github.com/Open-Argon/Isotope/src/indexof"
	initpkg "github.com/Open-Argon/Isotope/src/init"
	zipPack "github.com/Open-Argon/Isotope/src/package/zip"
)

const usage = `install [options] [package]`

var o = help.Options{
	{"specify a specific remote host", "--remote [host]"},
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

func installAddDependencies(dependencies []zipPack.Dependency, installing []zipPack.Dependency, remote string, path string, pkg zipPack.Dependency) {
	for _, dependency := range dependencies {
		for _, installingDependency := range installing {
			if dependency.Name == installingDependency.Name && dependency.Version == installingDependency.Version {
				log.Fatal("Circular dependency detected ", dependency.Name, dependency.Version)
			}
		}
		installPackage(remote, dependency.URL, dependency.Name, dependency.Version, path, append(installing, pkg))
	}
}

func installPackage(remote string, URL string, name string, version string, path string, installing []zipPack.Dependency) zipPack.Dependency {
	var pkg = zipPack.Dependency{
		Name:    name,
		Version: version,
		URL:     URL,
	}
	if URL == "" {
		params := url.Values{}
		params.Add("name", name)
		params.Add("version", version)
		URL = "https://" + remote + "/isotope-search?" + params.Encode()
		fmt.Println("Searching for", name+"@"+version, "(URL:", URL+")")
		resp, err := http.Get(URL)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Fatal("Package not found: ", resp.StatusCode)
		}
		pkg = zipPack.Dependency{}
		err = json.NewDecoder(resp.Body).Decode(&pkg)
		if err != nil {
			log.Fatal(err)
		}
		URL = pkg.URL
		fmt.Println("Found", pkg.Name+"@"+pkg.Version, "(URL:", URL+")")
		fmt.Println()
	}
	fmt.Println("Downloading from", URL)
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Package not found: ", resp.StatusCode)
	}
	tempFile, err := os.CreateTemp("", "isotope-download-*.zip")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	tempFile.Close()
	zipOpen, err := os.Open(tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	zipReader, err := gzip.NewReader(zipOpen)
	if err != nil {
		log.Fatal(err)
	}
	defer zipReader.Close()
	tarReader := tar.NewReader(zipReader)
	var dependencies = make([]zipPack.Dependency, 0)
	for {
		file, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		switch file.Name {
		case "iso-lock.json":
			err = json.NewDecoder(tarReader).Decode(&dependencies)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if pkg.Name == "" {
		log.Fatal("Package not valid")
	}
	argon_modules := filepath.Join(path, "argon_modules")
	modulepath := filepath.Join(argon_modules, pkg.Name)
	tempDir, err := os.MkdirTemp("", "isotope-package-install-"+pkg.Name+"-"+pkg.Version+"-")
	if err != nil {
		log.Fatal(err)
	}
	zipOpen.Seek(0, io.SeekStart)
	zipReader, err = gzip.NewReader(zipOpen)
	if err != nil {
		log.Fatal(err)
	}
	tarReader = tar.NewReader(zipReader)
	for {
		file, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		p := filepath.Join(tempDir, file.Name)
		os.MkdirAll(filepath.Dir(p), os.ModePerm)
		if !file.FileInfo().IsDir() {
			fileWriter, err := os.Create(p)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(fileWriter, tarReader)
			if err != nil {
				log.Fatal(err)
			}
			fileWriter.Close()
		}
	}
	installAddDependencies(dependencies, installing, remote, tempDir, pkg)
	deleteDirectoryIfExists(modulepath)
	os.MkdirAll(argon_modules, os.ModePerm)
	err = os.Rename(tempDir, modulepath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Installed", pkg.Name+"@"+pkg.Version)
	return pkg
}

func Install() {
	var args = args.Args[1:]
	if len(args) == 0 {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		lockPath := filepath.Join(path, "iso-lock.json")
		lockFile, err := os.Open(lockPath)
		if err != nil {
			log.Fatal(err)
		}
		var dependencies = make([]zipPack.Dependency, 0)
		err = json.NewDecoder(lockFile).Decode(&dependencies)
		if err != nil {
			log.Fatal(err)
		}
		lockFile.Close()
		for _, dependency := range dependencies {
			installPackage("", dependency.URL, dependency.Name, dependency.Version, path, make([]zipPack.Dependency, 0))
		}
		return
	}

	global := false
	if indexof.Indexof(args[1:], "--global") != -1 || indexof.Indexof(args[1:], "-g") != -1 {
		global = true
	}
	name := args[0]
	remote := config.DefaultRemote
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
	split := strings.SplitN(name, "@", 2)
	if len(split) == 2 {
		name = split[0]
		version = split[1]
	}
	var path string
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
	pkg := installPackage(remote, "", name, version, path, make([]zipPack.Dependency, 0))
	if !global {
		initpkg.Init(true, path)
		lockPath := filepath.Join(path, "iso-lock.json")
		lockFile, err := os.Open(lockPath)
		if err != nil {
			log.Fatal(err)
		}
		var dependencies = make([]zipPack.Dependency, 0)
		err = json.NewDecoder(lockFile).Decode(&dependencies)
		if err != nil {
			log.Fatal(err)
		}
		lockFile.Close()
		var newDependencies = make([]zipPack.Dependency, 0)
		for _, dependency := range dependencies {
			if dependency.Name != pkg.Name {
				newDependencies = append(newDependencies, dependency)
			}
		}
		newDependencies = append(newDependencies, pkg)
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

	fmt.Println("Install path:", path)
	fmt.Println()
	fmt.Println("Done!")
}
