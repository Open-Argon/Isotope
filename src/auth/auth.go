package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Open-Argon/Isotope/src/config"
	"github.com/Open-Argon/Isotope/src/help"
)

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

var usage = `auth [options]`
var o = help.Options{
	{"show help", "--help, -h"},
	{"add a new auth", "--add [host] [auth]"},
	{"remove an auth", "--remove [host]"},
}

func Auth() {
	dir := filepath.Dir(config.GlobalPath)
	authPath := filepath.Join(dir, "isotope-auth.json")
	auths := make(map[string]string)
	if fileExists(authPath) {
		reader, err := os.Open(authPath)
		if err != nil {
			log.Fatal(err)
		}
		err = json.NewDecoder(reader).Decode(&auths)
		if err != nil {
			log.Fatal(err)
		}
		reader.Close()
	}
	if len(os.Args) > 1 {
		switch os.Args[2] {
		case "--help", "-h":
			help.Help(usage, o)
			return
		case "--add":
			if len(os.Args) < 4 {
				help.Help(usage, o)
				return
			}
			auths[os.Args[3]] = os.Args[4]
			writer, err := os.Create(authPath)
			if err != nil {
				log.Fatal(err)
			}
			err = json.NewEncoder(writer).Encode(auths)
			if err != nil {
				log.Fatal(err)
			}
			writer.Close()
			fmt.Println("Added auth for", os.Args[3])
			return
		case "--remove":
			if len(os.Args) < 3 {
				help.Help(usage, o)
				return
			}
			if _, ok := auths[os.Args[3]]; !ok {
				fmt.Println("Auth for", os.Args[3], "not found")
				return
			}
			delete(auths, os.Args[3])
			writer, err := os.Create(authPath)
			if err != nil {
				log.Fatal(err)
			}
			err = json.NewEncoder(writer).Encode(auths)
			if err != nil {
				log.Fatal(err)
			}
			writer.Close()
			fmt.Println("Removed auth for", os.Args[3])
			return
		}
	}
	help.Help(usage, o)
}

func GetAuth(domain string) (string, error) {
	dir := filepath.Dir(config.GlobalPath)
	authPath := filepath.Join(dir, "isotope-auth.json")
	auths := make(map[string]string)
	if fileExists(authPath) {
		reader, err := os.Open(authPath)
		if err != nil {
			return "", err
		}
		err = json.NewDecoder(reader).Decode(&auths)
		if err != nil {
			return "", err
		}
		reader.Close()
	}
	return auths[domain], nil
}
