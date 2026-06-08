package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	Dir    string
	Header *multipart.FileHeader
	File   *multipart.File
}

type fileKind string

const (
	FileKindInstrument fileKind = "instrument"
)

func constructKeyFromFileName(kind fileKind, id any, file FileInfo) (string, error) {
	nameSplits := strings.Split(file.Header.Filename, ".")
	splitCount := len(nameSplits)
	if splitCount < 2 {
		return "", fmt.Errorf("invalid file name: %s\n", file.Header.Filename)
	}
	ext := nameSplits[splitCount-1]

	return filepath.Join("instrument", fmt.Sprintf("%d.%s", id, ext)), nil
}
