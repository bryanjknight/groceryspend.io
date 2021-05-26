package parser

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
)

func getTestDataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(
		filepath.Dir(
			filepath.Dir(
				filepath.Dir(
					filepath.Dir(filename)))), "test", "data")
}

func readFileAsString(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(b)
}
