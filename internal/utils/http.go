package utils

import (
	"mime"
	"net/http"
)

func GetFileType(data []byte) string {
	// buffer size is standardized to 512 bytes because its enough to get the file extension
	bufferSize := min(len(data), 512)

	contentType := http.DetectContentType(data[:bufferSize])

	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil || len(extensions) == 0 {
		// default to bin
		return ".bin"
	}

	return extensions[0]
}
