package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func FileSum(path string) string {
	content, err := os.ReadFile(path)
	Must(err)

	sum := md5.Sum([]byte(content))
	return hex.EncodeToString(sum[:])
}

func WriteFile(path, content string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	Must(err)
	defer f.Close()

	_, err = f.Write([]byte(content))
	Must(err)
}

func main() {
	_, err := exec.LookPath("wire")
	Must(err)

	for _, relativePath := range os.Args[1:] {
		wd, err := os.Getwd()
		Must(err)

		absPath := path.Join(wd, relativePath)
		sumPath := path.Join(wd, "./tmp/cached", relativePath, "sum")
		wirePath := path.Join(absPath, "wire.go")

		fmt.Printf("wiring %s... ", relativePath)

		_, err = os.Stat(sumPath)
		if !errors.Is(err, os.ErrNotExist) && err != nil {
			Must(err)
		} else if !errors.Is(err, os.ErrNotExist) {
			bs, err := os.ReadFile(sumPath)
			Must(err)

			if string(bs) == FileSum(wirePath) {
				fmt.Println("CACHED")
				continue
			}
		}

		fmt.Println()

		cmd := exec.Command("wire", absPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		Must(cmd.Run())

		Must(os.MkdirAll(path.Join(wd, "./tmp/cached", relativePath), 0700))

		WriteFile(sumPath, FileSum(wirePath))
	}
}
