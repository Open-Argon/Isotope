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
)

type Package map[string]any

func ReadPackageAndDependencies(path string) (Package, *bytes.Buffer) {
	packageFile, err := os.ReadFile(filepath.Join(path, "argon-package.json"))
	if err != nil {
		log.Fatal("failed to read argon-package.json")
	}

	var pkg Package
	if err := json.Unmarshal(packageFile, &pkg); err != nil {
		log.Fatal("failed to unmarshal argon-package.json")
	}

	if _, ok := pkg["name"]; !ok {
		log.Fatal("package name not found")
	} else if _, ok := pkg["version"]; !ok {
		log.Fatal("package version not found")
	} else if _, ok := pkg["dependencies"]; !ok {
		log.Fatal("package dependencies list not found")
	} else if _, ok := pkg["name"].(string); !ok {
		log.Fatal("package name is not a string")
	} else if _, ok := pkg["version"].(string); !ok {
		log.Fatal("package name is not a string")
	} else if _, ok := pkg["dependencies"].([]interface{}); !ok {
		log.Fatal("package dependencies list is not an array")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("src directory does not exist")
	}
	buf := new(bytes.Buffer)

	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
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

	return pkg, buf
}
