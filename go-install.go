package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	versionURL = "https://go.dev/VERSION?m=text"

	// https://dl.google.com/go/go1.16.3.linux-amd64.tar.gz
	downloadURL = "https://dl.google.com/go/%s.%s-%s.tar.gz"

	installPath = "/usr/local"
)

var (
	downloadFile = filepath.Join(os.TempDir(), "%s.%s-%s.tar.gz")
)

func main() {
	newestVersion, err := getNewestVersion()
	if err != nil {
		log.Fatal(err)
	}

	if newestVersion <= runtime.Version() {
		fmt.Println("already up to date")
		return
	}

	fmt.Println("updating from", runtime.Version(), "to", newestVersion)

	src, err := downloadTAR(newestVersion)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(src)

	err = extract(src)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done")
}

func getNewestVersion() (string, error) {
	res, err := http.Get(versionURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
