package experiment

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Experiment struct {
	File      string `yaml:"file"`
	Operation string `yaml:"operation"`
	Key       string `yaml:"key"`
}


func(e *Experiment) Execute() error {
	ops := strings.Split(e.Operation, " ")
	cmd := exec.Command(ops[0], ops[1], e.File)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	return nil
}

func(e *Experiment) Backup() error {
	return copyNative(e.File, "./backups/")
}

func(e *Experiment) RestoreFile() error {
	parts := strings.Split(e.File, "/")
	file := parts[len(parts) - 1]
	l := fmt.Sprintf("./backups/%s", file)
	return copyNative(l, e.File)
}
func copyNative(src, dst string) error {
	cmd := exec.Command("cp", "-rp", src, dst)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	return nil
}