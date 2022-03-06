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

	f, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer f.Close()

	res, err := http.Get(url)
	if err != nil {
		return path, err
	}
	defer res.Body.Close()

	bar := progressbar.DefaultBytes(
		res.ContentLength,
		"downloading",
	)

	_, err = io.Copy(io.MultiWriter(f, bar), res.Body)
	return path, err
}
