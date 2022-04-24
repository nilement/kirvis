package backups

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func RestoreFile(origFile string) error {
	parts := strings.Split(origFile, "/")
	file := parts[len(parts) - 1]
	backup := fmt.Sprintf("./backups/%s", file)

	return copyFile(backup, origFile)
}

func BackupFile(apiFile string) error {
	parts := strings.Split(apiFile, "/")
	file := parts[len(parts) - 1]
	backup := fmt.Sprintf("./backups/%s", file)

	return copyFile(apiFile, backup)
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