//go:build linux && amd64

package tools

import "fmt"

func GetUrl(version string) string {
	return fmt.Sprintf("https://go.dev/dl/go%v.linux-amd64.tar.gz", version)
}
