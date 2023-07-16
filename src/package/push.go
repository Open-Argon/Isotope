package pack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Open-Argon/Isotope/src/hash"
)

func push() {
	var name string
	var version string
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	packageFilePath := filepath.Join(path, "iso-package.json")
	packageFile, err := os.ReadFile(packageFilePath)
	if err != nil {
		log.Fatal(err)
	}
	var pkg map[string]any
	if err := json.Unmarshal(packageFile, &pkg); err != nil {
		log.Fatal(err)
	}
	if nameAny, ok := pkg["name"]; ok {
		name = nameAny.(string)
	} else {
		log.Fatal("package name not found")
	}
	if versionAny, ok := pkg["version"]; ok {
		version = versionAny.(string)
	} else {
		log.Fatal("package version not found")
	}
	hash := hash.Sha256Hex(name + "@" + version)
	zipPath := filepath.Join(path, "__isotope__", "builds", "armod-"+hash+".tar.gz")
	zipFile, err := os.Open(zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()
	fmt.Println("Pushing", name+"@"+version)
	fmt.Println("Path:", zipPath)
	fmt.Println("Hash:", hash)

	// make a http post request to the server
	req, err := http.NewRequest("POST", "https://pkg.argon.wbell.dev/push", zipFile)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Package-Name", name)
	req.Header.Set("Package-Version", version)
	req.Header.Set("Package-Hash", hash)

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Done!")
}
