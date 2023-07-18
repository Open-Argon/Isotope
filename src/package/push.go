package pack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Open-Argon/Isotope/src/args"
	"github.com/Open-Argon/Isotope/src/config"
	"github.com/Open-Argon/Isotope/src/hash"
)

func push() {
	var name string
	var version string
	var remote = config.DefaultRemote
	var latest = false
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	packageFilePath := filepath.Join(path, "argon-package.json")
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
	if len(args.Args) > 2 {
		if args.Args[2] == "--latest" {
			latest = true
			if len(args.Args) > 3 {
				remote = args.Args[3]
			}
		} else {
			remote = args.Args[2]
		}
	}
	hash := hash.Sha256Hex(name + "@" + version)
	zipPath := filepath.Join(path, "__isotope__", "builds", "armod-"+hash+".tar.gz")
	zipFile, err := os.Open(zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()
	fmt.Println("Pushing", name+"@"+version)
	fmt.Println("Remote:", remote)
	fmt.Println("Path:", zipPath)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "file.tar.gz")
	io.Copy(part, zipFile)
	writer.Close()
	writer.WriteField("name", name)
	writer.WriteField("version", version)
	writer.WriteField("latest", fmt.Sprint(latest))
	r, err := http.NewRequest("POST", "https://"+remote+"/isotope-push", body)
	if err != nil {
		log.Fatal(err)
	}
	r.Method = "POST"
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal("Package not found: ", resp.StatusCode)
	}
	fmt.Println("Pushed", name+"@"+version)
}
