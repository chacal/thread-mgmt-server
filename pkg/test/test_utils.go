package test

import (
	"io/ioutil"
	"os"
)

func Tempfile() string {
	f, err := ioutil.TempFile("", "devices-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
