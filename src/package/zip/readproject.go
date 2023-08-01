package zip

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	validpackagename "github.com/Open-Argon/Isotope/src/validPackageName"
)

type Package struct {
	Name         string
	Version      string
	Dependencies []Dependency
}

type Dependency struct {
	Name    string
	Version string
	URL     string
	Remote  string
}

func ReadPackageAndDependencies(path string) (Package, *bytes.Buffer) {
	src := filepath.Join(path, "src")
	packageFilePath := filepath.Join(path, "argon-package.json")
	packageFile, err := os.ReadFile(packageFilePath)
	if err != nil {
		log.Fatal(err)
	}
	LockFilePath := filepath.Join(path, "iso-lock.json")
	LockFile, err := os.ReadFile(LockFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var pkg map[string]any
	if err := json.Unmarshal(packageFile, &pkg); err != nil {
		log.Fatal(err)
	}

	var pkgObj Package

	var lock []Dependency
	if err := json.Unmarshal(LockFile, &lock); err != nil {
		log.Fatal(err)
	}
	name, ok := pkg["name"]
	if !ok && name.(string) == "" {
		log.Fatal("package name not found")
	}
	if validpackagename.ValidPackageName(name.(string)) != name.(string) {
		log.Fatal("package name is invalid")
	}
	version, ok := pkg["version"]
	if !ok && version.(string) == "" {
		log.Fatal("package version not found")
	}

	pkgObj.Name = name.(string)
	pkgObj.Version = version.(string)
	pkgObj.Dependencies = lock

	if _, err := os.Stat(src); os.IsNotExist(err) {
		log.Fatal("src directory not found")
	}

	buf := new(bytes.Buffer)

	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	err = addToArchive(tw, packageFilePath, "argon-package.json")
	if err != nil {
		log.Fatal(err)
	}
	err = addToArchive(tw, LockFilePath, "iso-lock.json")
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(src, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		saveAs, err := filepath.Rel(src, filePath)
		if err != nil {
			return err
		}
		err = addToArchive(tw, filePath, saveAs)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("failed to create zip file contents: %w", err))
	}

	return pkgObj, buf
}

func addToArchive(tw *tar.Writer, path string, saveAs string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = saveAs

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}
