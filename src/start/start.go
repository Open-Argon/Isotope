package start

import (
	"os"
	"path/filepath"
)

var ZippedPath = ""

func init() {
	exc, err := os.Executable()
	if err != nil {
		panic(err)
	}
	excPath := filepath.Dir(exc)
	ZippedPath = filepath.Join(excPath, "zipped")
	err = os.MkdirAll(ZippedPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

}
