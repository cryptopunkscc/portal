package assets

import (
	"io"
	"io/fs"
	"net/http"
	"strings"
)

const (
	TypeDir  = "fs/dir"
	TypeHtml = "text/html"
	TypeJs   = "text/js"
	TypeZip  = "application/zip"
)

func FileType(files fs.FS, filePath string) (string, error) {
	file, err := files.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 10512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := strings.Split(http.DetectContentType(buffer), ";")[0]

	return contentType, nil
}
