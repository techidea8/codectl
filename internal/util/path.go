package util

import (
	"os"
	"path/filepath"
)

func CurrentExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	if exePath, err = filepath.Abs(exePath); err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exePath)
}

func CurrentExecutableDir() (string, error) {
	if f, e := CurrentExecutablePath(); e != nil {
		return "", e
	} else {
		return filepath.Dir(f), nil
	}

}
