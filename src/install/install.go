package install

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/help"
	"github.com/Open-Argon/Isotope/src/indexof"
	zipPack "github.com/Open-Argon/Isotope/src/package/zip"
)

var usage = `install [package] [options]`
var o = help.Options{
	{"specify a specific remote host", "--remote [host]"},
	{"show help", "--help, -h"},
}

func deleteFilesAndDirectories(path string) error {
	err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		err = os.Remove(filePath)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
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
		URL = "http://" + remote + "/isotope-search?" + params.Encode()
		fmt.Println("Searching for", name+"@"+version, "(URL:", URL+")")
		resp, err := http.Get(URL)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Fatal("Package not found")
		}
		pkg = zipPack.Dependency{}
		err = json.NewDecoder(resp.Body).Decode(&pkg)
		if err != nil {
			log.Fatal(err)
		}
		URL = pkg.URL
	}
	fmt.Println("Downloading from", URL)
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Package not found")
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
	modulepath := filepath.Join(path, "argon_modules", pkg.Name)
	tempDir, err := ioutil.TempDir("", "mytempdir")
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
	for _, dependency := range dependencies {
		for _, installingDependency := range installing {
			if dependency.Name == installingDependency.Name && dependency.Version == installingDependency.Version {
				log.Fatal("Circular dependency detected ", dependency.Name, dependency.Version)
			}
		}
		installPackage(remote, dependency.URL, dependency.Name, dependency.Version, tempDir, append(installing, pkg))
	}
	deleteFilesAndDirectories(modulepath)
	os.MkdirAll(modulepath, os.ModePerm)
	err = os.Rename(tempDir, modulepath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Installed", pkg.Name+"@"+pkg.Version)
	fmt.Println("Path:", modulepath)
	return pkg
}

func Install() {
	var args = args.Args[1:]
	if len(args) == 0 {
		help.Help(usage, o)
		return
	}

	global := false
	if indexof.Indexof(args[1:], "--global") != -1 || indexof.Indexof(args[1:], "-g") != -1 {
		global = true
	}
	name := args[0]
	remote := "localhost:3000"
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
	installPackage(remote, "", name, version, path, make([]zipPack.Dependency, 0))
}
