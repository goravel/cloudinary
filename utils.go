package cloudinary

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

// GetRawContent retrieves the raw content of a file from the provided URL.
func GetRawContent(url string) ([]byte, error) {
	// Make an HTTP GET request to fetch the file data
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching raw content: %w", err)
	}

	rawContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return rawContent, nil
}

func validPath(path string) string {
	realPath := strings.TrimPrefix(path, "."+string(filepath.Separator))
	realPath = strings.TrimPrefix(realPath, string(filepath.Separator))
	realPath = strings.TrimPrefix(realPath, ".")
	realPath = strings.TrimSuffix(realPath, string(filepath.Separator))
	return realPath
}
