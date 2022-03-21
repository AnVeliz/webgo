package fileutils

import (
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
)

func Download(rootDir, urlStr string) error {
	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	urlValue, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	fileName := path.Join(rootDir, urlValue.Path)
	out, err := CreateLocal(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func CreateLocal(p string) (*os.File, error) {
	dirPath := path.Dir(p)
	if err := os.MkdirAll(dirPath, fs.ModeDir); err != nil {
		return nil, err
	}
	return os.Create(p)
}
