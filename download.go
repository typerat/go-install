package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type Progress struct {
	Start   time.Time
	Current uint64
	Total   uint64
}

func (p *Progress) Write(b []byte) (int, error) {
	n := len(b)
	p.Current += uint64(n)

	rate := p.Current * 1e9 / uint64(time.Since(p.Start))
	percent := p.Current * 100 / p.Total

	fmt.Printf("\r%s", strings.Repeat(" ", 100))
	fmt.Printf("\r%d %%\t%s\t %s/s", percent, humanize.Bytes(p.Current), humanize.Bytes(rate))

	return n, nil
}

func NewProgress(size uint64) *Progress {
	return &Progress{
		Start: time.Now(),
		Total: size,
	}
}

func downloadTAR(version string) (string, error) {
	url := fmt.Sprintf(downloadURL, version, runtime.GOOS, runtime.GOARCH)
	path := fmt.Sprintf(downloadFile, version, runtime.GOOS, runtime.GOARCH)

	fmt.Println("downloading", url, "to", path)

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

	size, err := strconv.ParseUint(res.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return path, err
	}

	progress := NewProgress(size)
	_, err = io.Copy(f, io.TeeReader(res.Body, progress))
	return path, err
}
