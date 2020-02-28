package http

import (
	"fmt"
	"net/http"
	"os"
	"pet-spotlight/io"
)

// Download downloads the file from the specified URL and saves to the provided path as the specified file
// name.
func Download(url string, path string, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get image from %s: %w", url, err)
	}
	defer io.CloseResource(resp.Body)
	filePath := fmt.Sprintf("%s/%s", path, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer io.CloseResource(f)
	if err = io.CopyToFile(resp.Body, f); err != nil {
		return err
	}
	return nil
}
