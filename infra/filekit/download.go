package filekit

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

func DownloadAs(appdir, url string) (filepath string, err error) {
	fileNameWithSuffix := path.Base(url)
	fileNameWithSuffix = strings.Split(fileNameWithSuffix, "?")[0]
	fileType := path.Ext(fileNameWithSuffix)
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	uuidValue := uuid.New()
	filepath = path.Join(appdir, uuidValue.String()+fileType)
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(filepath, bytes, os.ModeAppend)
	if err != nil {
		return "", err
	}
	return filepath, nil
}
