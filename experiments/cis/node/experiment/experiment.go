package experiment

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	parts := strings.Split(e.File, "/")
	file := parts[len(parts) - 1]
	backup := fmt.Sprintf("./backups/%s", file)

	return copyFile(e.File, backup)
}

func(e *Experiment) RestoreFile() error {
	parts := strings.Split(e.File, "/")
	file := parts[len(parts) - 1]
	l := fmt.Sprintf("./backups/%s", file)
	return copyFile(l, e.File)
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}