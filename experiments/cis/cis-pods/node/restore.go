package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

func restoreFiles() error {
	s, err := readFiles("./restore.yaml")
	if err != nil {
		log.Fatal(err)
	}
	for _, str := range *s {
		n := strings.Split(str, "/")
		err = copyNative("./backups/"+n[len(n)-1], str)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func readFiles(filename string) (*[]string, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c []string
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return &c, nil
}
