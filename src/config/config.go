package config

import (
	"fmt"
	"os"
	"path"
)

const DefaultRemote = "isotope.wbell.dev"

var GlobalPath string

func init() {
	var home_dir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
	}
	GlobalPath = path.Join(home_dir, ".argon")
}
