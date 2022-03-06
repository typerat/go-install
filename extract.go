package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func extract(src string) error {
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

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	// add progress bar
	bar := progressbar.DefaultBytes(
		stat.Size(),
		"extracting",
	)
	r := io.TeeReader(f, bar)

	plain, err := gzip.NewReader(r)
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
