package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/schollz/progressbar/v3"
)

func downloadTAR(version string) (string, error) {
	url := fmt.Sprintf(downloadURL, version, runtime.GOOS, runtime.GOARCH)
	path := fmt.Sprintf(downloadFile, version, runtime.GOOS, runtime.GOARCH)

	err := os.RemoveAll(path)
	if err != nil {
		return "", err
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bar := progressbar.DefaultBytes(
		res.ContentLength,
		"downloading",
	)

	_, err = io.Copy(io.MultiWriter(f, bar), res.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}
