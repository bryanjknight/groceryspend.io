package utils

import "io/ioutil"

// ReadFileAsString reads a local file and returns it as a string
func ReadFileAsString(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(b)
}
