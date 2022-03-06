package main

import (
	"archive/tar"
	"compress/gzip"
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

	err = install(src)
	if err != nil {
		log.Fatal(err)
	}
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

func install(src string) error {
	fmt.Println("extracting", src, "to", installPath)

	goroot := filepath.Join(installPath, "go")
	err := os.RemoveAll(goroot)
	if err != nil {
		return err
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	plain, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tr := tar.NewReader(plain)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		path := filepath.Join(installPath, header.Name)
		dir := filepath.Dir(path)

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		fmt.Println(path)
		out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
		if err != nil {
			return err
		}

		_, err = io.Copy(out, tr)
		if err != nil {
			return err
		}

		out.Close()
	}

	return nil
}
