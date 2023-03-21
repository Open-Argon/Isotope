package start

import (
	"os"
	"path"
)

var ZippedPath = ""

func init() {
	exc, err := os.Executable()
	if err != nil {
		panic(err)
	}
	ZippedPath = path.Join(path.Dir(exc), "__zipped__")
	err = os.MkdirAll(ZippedPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

}
