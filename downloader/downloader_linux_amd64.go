//go:build linux && amd64

package downloader

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

func Download(ver string) (*os.File, error) {
	//url := fmt.Sprintf("https://go.dev/dl/go%v.linux-amd64.tar.gz", ver)
	return nil, nil
}

func GetLatestVer() (string, error) {
	resp, err := http.Get("https://go.dev/VERSION?m=text")
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(http.StatusText(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	bodystring := string(body)
	versionnum := strings.Split(bodystring, "\n")[0][2:]
	return versionnum, nil
}
