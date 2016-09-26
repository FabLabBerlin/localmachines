package models

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func UploadImage(name string, dataUri string) (err error) {
	if err = checkFilename(name); err != nil {
		return
	}
	f, err := os.OpenFile("files/"+name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}

	data, err := decodeDataUri(dataUri)
	if err != nil {
		return fmt.Errorf("decode data uri: %v", err)
	}

	_, err = f.Write(data)
	defer f.Close()
	return
}

func decodeDataUri(dataUri string) ([]byte, error) {
	sep := "base64,"
	i := strings.Index(dataUri, sep)
	if i < 0 {
		return nil, fmt.Errorf("cannot remove prefix from data uri")
	}
	dataUri = dataUri[i+len(sep):]
	return base64.StdEncoding.DecodeString(dataUri)
}

func FileExtensionFor(mimeType string) (t string, err error) {
	switch mimeType {
	case "image/svg+xml":
		return ".svg", nil
	}
	return "", fmt.Errorf("cannot determine file extension for mime type '%s'", t)
}

func MimeType(name string) (t string, err error) {
	if strings.HasSuffix(name, ".svg") {
		return "image/svg+xml", nil
	}
	return "", fmt.Errorf("cannot determine mime type of %s", name)
}

// checkFilename, so it's not possible to do thing like
// checkFilename("../../../../etc/passwd")
func checkFilename(name string) (err error) {
	forbidden := []string{
		"..",
		"~",
		"`",
		"'",
		"\"",
		"!",
		"$",
		"%",
	}
	for _, f := range forbidden {
		if strings.Contains(name, f) {
			return fmt.Errorf("file name contains %s", f)
		}
	}
	return nil
}
