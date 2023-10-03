package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

	currentVersion, err := getCurrentVersion()
	if err != nil {
		log.Fatal(err)
	}

	if newestVersion <= currentVersion {
		fmt.Printf("already up to date (%s)\n", currentVersion)
		return
	}

	fmt.Println("updating from", currentVersion, "to", newestVersion)

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

func getCurrentVersion() (string, error) {
	buf, err := exec.Command(filepath.Join(installPath, "go", "bin", "go"), "version").Output()
	if err != nil {
		return "", nil
	}

	buf = regexp.MustCompile(`go\d\.\d+(\.\d+)?`).Find(buf)

	return string(buf), nil
}
