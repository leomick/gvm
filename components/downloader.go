package components

import (
	"net/http"
	"strconv"

	"github.com/charmbracelet/bubbles/progress"
)

type downloaderModel struct {
	bar            progress.Model
	percentage     uint8
	url            string
	downloadedSize uint64
	totalSize      uint64
}

func New(url string) (downloaderModel, error) {
	resp, err := http.Head(url)
	if err != nil {
		return downloaderModel{}, err
	}
	contentSize := resp.Header.Get("Content-Length")
	totalSizeInt, err := strconv.Atoi(contentSize)
	if err != nil {
		return downloaderModel{}, err
	}
	totalSize := uint64(totalSizeInt)
	bar := progress.New()
	return downloaderModel{
		bar:        bar,
		percentage: 0,
		url:        url,
		totalSize:  totalSize,
	}, nil
}
