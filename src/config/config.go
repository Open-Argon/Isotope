package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const DefaultRemote = "isotope.wbell.dev"

var GlobalPath string

func init() {
	var executable, err = os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
	}

	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		fmt.Println("Error:", err)
	}
	GlobalPath = filepath.Dir(filepath.Dir(executable))
}
