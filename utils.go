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

func defaultResourcesTypes() map[string][]string {
	return map[string][]string{
		"image": {"ai", "avif", "bmp", "bw", "djvu", "dng", "ps", "ept", "eps", "eps3", "fbx", "flif", "gif", "glb", "heif", "heic", "ico", "indd", "jpg", "jpe", "jpeg", "jp2", "wdp", "jxr", "hdp", "jxl", "obj", "pdf", "ply", "png", "psd", "arw", "cr2", "svg", "tga", "tif", "tiff", "u3ma", "usdz", "webp"},
		"video": {"3g2", "3gp", "avi", "flv", "m3u8", "ts", "m2ts", "mts", "mov", "mkv", "mp4", "mpeg", "mpd", "mxf", "ogv", "webm", "wmv", "aac", "aiff", "amr", "flac", "m4a", "mp3", "ogg", "opus", "wav"},
		"raw":   {"txt", "doc", "rtf", "xls", "ppt", "docx", "xlsx", "pptx", "pps", "ppsx", "odt", "ods", "odp"},
	}
}
