package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/AnVeliz/webgo/internal/chromium"
)

var (
	//go:embed assets/*
	embeddedContent embed.FS
)

func main() {
	baseUrl := "http://localhost"
	port := 8088
	address := fmt.Sprintf("%s:%d/", baseUrl, port)

	mux := http.NewServeMux()
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}

	mux.Handle("/", http.FileServer(http.FS(embeddedContent)))
	log.Println("main file server handler has set up")

	go func() {
		time.Sleep(time.Duration(1 * time.Second))
		chromiumTmpDir := createTemporaryChromium()
		defer os.RemoveAll(chromiumTmpDir)

		for _, file := range chromium.Files {
			downloadFile(chromiumTmpDir, fmt.Sprintf("%s%s", address, file))
		}

		appRootFile := fmt.Sprintf("%s%s", address, "assets/webui/index.html")
		cmd := exec.Command(path.Join(chromiumTmpDir, "assets/chromium/99.0.4844.74_x64/Chrome-bin/chrome.exe"), fmt.Sprintf("--app=%s", appRootFile))
		cmd.Run()
		srv.Shutdown(context.Background())
	}()

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	fmt.Println("main file server handler has stopped")
}

func createTemporaryChromium() string {
	chromiumDir, err := ioutil.TempDir("", "chromium")
	if err != nil {
		log.Fatal(err)
	}

	return chromiumDir
}

func downloadFile(rootDir, urlStr string) error {
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
	out, err := createFile(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func createFile(p string) (*os.File, error) {
	dirPath := path.Dir(p)
	if err := os.MkdirAll(dirPath, fs.ModeDir); err != nil {
		return nil, err
	}
	return os.Create(p)
}
