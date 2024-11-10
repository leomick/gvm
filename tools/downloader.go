package tools

import (
	"fmt"
	"io"
	"leomick/gvm/tools/targz"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func Download(ver string) error {
	url := GetUrl(ver)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 404 {
			return fmt.Errorf("404 %v. The version number you specified is probably invalid", http.StatusText(resp.StatusCode))
		}
		return fmt.Errorf("%v %v", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	err = targz.ExtractTarGZ(resp.Body, viper.GetString("installDir"), renamer(ver))
	if err != nil {
		return err
	}
	return nil
}

func GetLatestVer() (string, error) {
	resp, err := http.Get("https://go.dev/VERSION?m=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%v %v", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	bodystring := string(body)
	versionnum := strings.Split(bodystring, "\n")[0][2:]
	return versionnum, nil
}

func renamer(ver string) func(string) string {
	return func(name string) string {
		return strings.Replace(name, "go", ver, 1)
	}
}
