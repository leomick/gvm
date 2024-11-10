package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type progressWriter struct {
	total      int
	downloaded int
	Resp       *http.Response
	Content    []byte
}

type DoneMsg bool
type progressMsg float64

var ProgressChannel chan float64

func (pw *progressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	_, err := io.Copy(pw, pw.Resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	pw.Content = append(pw.Content, p...)
	if pw.total > 0 {
		pw.onProgress(float64(pw.downloaded) / float64(pw.total))
	}
	return len(p), nil
}

func (pw *progressWriter) onProgress(prog float64) {
	ProgressChannel <- prog
}

func WaitForProgress() tea.Msg {
	progress := <-ProgressChannel
	if progress == 1 {
		return DoneMsg(true)
	}
	return progressMsg(progress)
}

func (m Model) Start() tea.Cmd {
	go m.Pw.Start()
	return WaitForProgress
}

type Model struct {
	bar      progress.Model
	Pw       *progressWriter
	url      string
	progress float64
}

func New(url string) Model {
	ProgressChannel = make(chan float64)
	bar := progress.New()
	bar.FullColor = "123"
	bar.EmptyColor = "8"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal(fmt.Errorf("receiving status of %v from the url: %v", resp.Status, url))
	}
	if resp.ContentLength <= 0 {
		log.Fatal("Error when getting content length")
	}
	pw := &progressWriter{
		total: int(resp.ContentLength),
		Resp:  resp,
	}
	return Model{
		bar: bar,
		Pw:  pw,
		url: url,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case progressMsg:
		m.progress = float64(msg)
		cmd := WaitForProgress
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.bar.ViewAs(m.progress)
}
