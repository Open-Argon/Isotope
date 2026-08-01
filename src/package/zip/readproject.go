package zip

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	validpackagename "github.com/Open-Argon/Isotope/src/validPackageName"
	ignore "github.com/sabhiram/go-gitignore"
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

func loadIgnoreFile(containPath string) (*ignore.GitIgnore, error) {
	path := filepath.Join(containPath, ".isotopeignore")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ignore.CompileIgnoreLines(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	return ignore.CompileIgnoreLines(lines...), nil
}

func findLicenseFile(containPath string) string {
	candidates := []string{
		"LICENSE",
		"LICENSE.txt",
		"LICENSE.md",
		"LICENCE",
		"LICENCE.txt",
		"LICENCE.md",
	}

	for _, name := range candidates {
		path := filepath.Join(containPath, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func ReadPackageAndDependencies(contain_path string) (Package, *bytes.Buffer) {
	src := filepath.Join(contain_path, "src")
	packageFilePath := filepath.Join(contain_path, "argon-package.json")
	packageFile, err := os.ReadFile(packageFilePath)
	if err != nil {
		log.Fatal(err)
	}
	LockFilePath := filepath.Join(contain_path, "iso-lock.json")
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
	name, ok := pkg["name"].(string)
	if !ok || name == "" {
		log.Fatal("package name not found")
	}
	if validpackagename.ValidPackageName(name) != name {
		log.Fatal("package name is invalid")
	}
	version, ok := pkg["version"].(string)
	if !ok && version == "" {
		log.Fatal("package version not found")
	}

	pkgObj.Name = name
	pkgObj.Version = version
	pkgObj.Dependencies = lock

	if _, err := os.Stat(src); os.IsNotExist(err) {
		log.Fatal("src directory not found")
	}

	buf := new(bytes.Buffer)

	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	build, ok := pkg["build"].(string)
	if ok {
		buildFilePath := filepath.Join(contain_path, build)
		err = addToArchive(tw, buildFilePath, build)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = addToArchive(tw, packageFilePath, "argon-package.json")
	if err != nil {
		log.Fatal(err)
	}
	err = addToArchive(tw, LockFilePath, "iso-lock.json")
	if err != nil {
		log.Fatal(err)
	}

	licensePath := findLicenseFile(contain_path)
	if licensePath != "" {
		err = addToArchive(tw, licensePath, filepath.Base(licensePath))
		if err != nil {
			log.Fatal(err)
		}
	}

	isotopeIgnore, err := loadIgnoreFile(contain_path)
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(contain_path, path)
		if err != nil {
			return err
		}

		// gitignore expects forward slashes
		relPath = filepath.ToSlash(relPath)

		if isotopeIgnore.MatchesPath(relPath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() && d.Name() == "__arcache__" {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		return addToArchive(tw, path, relPath)
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
