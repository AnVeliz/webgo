package chromium

import (
	"embed"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/AnVeliz/webgo/internal/fileutils"
)

func Run(url string, fs embed.FS) {
	runChromium(url, fs)
}

func runChromium(url string, fs embed.FS) error {
	cmd, chromiumTmpDir := prepareChromiumCmd(url, fs)
	defer os.RemoveAll(chromiumTmpDir)

	if !checkConnection(url, time.Duration(1*time.Second), 5) {
		return errors.New("connection can not be established")
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func checkConnection(url string, timeout time.Duration, maxAttemptsNum int) bool {
	httpClient := http.Client{
		Timeout: timeout,
	}

	index := 0
	for resp, err := httpClient.Get(url); err != nil || resp.StatusCode != http.StatusOK; index++ {
		if index == maxAttemptsNum-1 {
			fmt.Printf("can not establish connection error: %s\n", err)
			return false
		}
		time.Sleep(timeout)
	}
	fmt.Printf("connected after %d attempt\n", index+1)

	return true
}

func prepareChromiumCmd(address string, fs embed.FS) (*exec.Cmd, string) {
	chromiumTmpDir := createTemporaryChromium()

	for _, file := range Files {
		fileutils.Save(chromiumTmpDir, file, fs)
	}

	appRootFile := fmt.Sprintf("%s%s", address, "index.html")
	cmd := exec.Command(path.Join(chromiumTmpDir, "assets/chromium/99.0.4844.74_x64/Chrome-bin/chrome.exe"), fmt.Sprintf("--app=%s", appRootFile))
	return cmd, chromiumTmpDir
}

func createTemporaryChromium() string {
	chromiumDir, err := ioutil.TempDir("", "chromium")
	if err != nil {
		log.Fatal(err)
	}

	return chromiumDir
}
