package io_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"pet-spotlight/io"
	"testing"
)

func TestWriteFile(t *testing.T) {
	f, err := ioutil.TempFile("", "write-file-test.txt")
	if err != nil {
		t.Fatal(err)
	}
	fileName := f.Name()
	defer os.Remove(fileName)
	defer f.Close()
	expectedContent := "this is a test"
	if err = io.WriteFile(expectedContent, fileName); err != nil {
		t.Error(err)
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Error(err)
	}
	if string(content) != expectedContent {
		t.Errorf("content of the file is %s", content)
	}
}

func TestCopyToFile(t *testing.T) {
	var content bytes.Buffer
	expectedContent := "this is a test"
	if _, err := content.WriteString(expectedContent); err != nil {
		t.Fatal(err)
	}
	var r bytes.Buffer
	if err := io.CopyToFile(&content, &r); err != nil {
		t.Fatal(err)
	}
	output := r.String()
	if output != expectedContent {
		t.Errorf("content of file is %s", output)
	}
}

type closerTester struct {
}

func (c closerTester) Close() error {
	return errors.New("failed")
}

func TestCloseResource(t *testing.T) {
	tester := closerTester{}
	io.CloseResource(tester)
}
