package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/andybalholm/crlf"
	"golang.org/x/text/transform"
)

const testdir = "testdata"

func readFile(name string) ([]byte, error) {
	golden, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	goldenstr := filepath.FromSlash(string(golden))
	goldenlf, _, err := transform.Bytes(new(crlf.Normalize), []byte(goldenstr))
	if err != nil {
		return nil, err
	}
	return goldenlf, nil
}

func yamlData(v interface{}) []byte {
	b, err := marshalYAML(v)
	if err != nil {
		panic(err)
	}
	return b
}
