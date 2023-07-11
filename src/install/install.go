package install

import (
	"archive/zip"
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
	"github.com/Open-Argon/Isotope/src/help"
	zipPack "github.com/Open-Argon/Isotope/src/package/zip"
)

var usage = `install [options] [package]`
var o = help.Options{
	{"specify a specific remote host", "--remote [host]"},
	{"show help", "--help, -h"},
}

var installing = make(map[string]bool)

func deleteFilesAndDirectories(path string) error {
	err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			err = os.Remove(filePath)
			if err != nil {
				return err
			}
		} else {
			err = os.Remove(filePath)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func installPackage(remote string, name string, version string, global bool) zipPack.Dependency {
	params := url.Values{}
	params.Add("name", name)
	params.Add("version", version)
	urlPath := "http://" + remote + "/download?" + params.Encode()
	fmt.Println("Searching for", name+"@"+version, "(url:", urlPath+")")
	resp, err := http.Get(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	tempFile, err := os.CreateTemp("", "download-*.zip")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	tempFile.Close()
	zipReader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer zipReader.Close()
	var pkg = zipPack.Dependency{}
	var dependencies = make([]zipPack.Dependency, 0)
	var path = ""
	if global {
		path, err = os.Executable()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		path, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, file := range zipReader.File {
		switch file.Name {
		case "iso-package.json":
			fileReader, err := file.Open()
			if err != nil {
				log.Fatal(err)
			}
			packageDecoded := make(map[string]any)
			err = json.NewDecoder(fileReader).Decode(&packageDecoded)
			if err != nil {
				log.Fatal(err)
			}
			pkg.Name = packageDecoded["name"].(string)
			pkg.Version = packageDecoded["version"].(string)
			params.Add("name", pkg.Name)
			params.Add("version", pkg.Version)
			pkg.URL = "http://" + remote + "/download?" + params.Encode()
			fileReader.Close()
		case "iso-package-lock.json":
			fileReader, err := file.Open()
			if err != nil {
				log.Fatal(err)
			}
			err = json.NewDecoder(fileReader).Decode(&dependencies)
			if err != nil {
				log.Fatal(err)
			}
			fileReader.Close()
		}
	}
	modulepath := filepath.Join(path, "argon_modules", pkg.Name)
	deleteFilesAndDirectories(modulepath)
	os.MkdirAll(modulepath, os.ModePerm)
	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		p := filepath.Join(modulepath, file.Name)
		os.MkdirAll(filepath.Dir(p), os.ModePerm)
		if !file.FileInfo().IsDir() {
			fileWriter, err := os.Create(p)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(fileWriter, fileReader)
			if err != nil {
				log.Fatal(err)
			}
			fileWriter.Close()
		}
	}
	for _, dependency := range dependencies {
		if !installing[dependency.Name] {
			installing[dependency.Name] = true
			installPackage(remote, dependency.Name, dependency.Version, global)
		}
	}
	return pkg
}

func Install() {
	var args = args.Args[1:]
	if len(args) == 0 {
		help.Help(usage, o)
		return
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
	installPackage(remote, name, version, false)
}
