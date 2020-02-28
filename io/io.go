package io

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// MakeDir creates the specified directory.
func MakeDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// WriteFile writes the string content to the specified file.
func WriteFile(content string, file string) error {
	if err := ioutil.WriteFile(file, []byte(content), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write content to %s: %w", file, err)
	}
	return nil
}

// CopyFile copies the content to the file.
func CopyToFile(content io.Reader, w io.Writer) error {
	if _, err := io.Copy(w, content); err != nil {
		return fmt.Errorf("failed to copy content to file: %w", err)
	}
	return nil
}

// CloseResource closes the provided closer.
func CloseResource(c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Println(err)
	}
}
