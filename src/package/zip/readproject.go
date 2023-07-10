package zip

import (
	"archive/zip"
	"bytes"
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
	Dependencies []dependency
}

type dependency struct {
	Name    string
	Version string
	URL     string
}

func ReadPackageAndDependencies(path string) (Package, *bytes.Buffer) {
	src := filepath.Join(path, "src")
	packageFile, err := os.ReadFile(filepath.Join(path, "iso-package.json"))
	if err != nil {
		log.Fatal(err)
	}
	LockFile, err := os.ReadFile(filepath.Join(path, "iso-package-lock.json"))
	if err != nil {
		log.Fatal(err)
	}

	var pkg map[string]any
	if err := json.Unmarshal(packageFile, &pkg); err != nil {
		log.Fatal(err)
	}

	var pkgObj Package

	var lock []dependency
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

	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()
	packageFileWriter, err := zipWriter.Create("iso-package.json")
	if err != nil {
		panic(fmt.Errorf("failed to create zip file contents: %w", err))
	}
	_, err = packageFileWriter.Write(packageFile)
	if err != nil {
		panic(fmt.Errorf("failed to create zip file contents: %w", err))
	}
	LockFileWriter, err := zipWriter.Create("iso-package-lock.json")
	if err != nil {
		panic(fmt.Errorf("failed to create zip file contents: %w", err))
	}
	_, err = LockFileWriter.Write(LockFile)
	if err != nil {
		panic(fmt.Errorf("failed to create zip file contents: %w", err))
	}
	err = filepath.Walk(src, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println("packaging", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		f, err := zipWriter.Create(filePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
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
