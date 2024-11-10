package downloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/codeclysm/extract/v4"
	"github.com/spf13/viper"
)

type progressWriter struct {
	total      int
	downloaded int
	resp       *http.Response
}

type doneMsg bool
type percentageMsg int64

var percentChannel chan float64

func (pw *progressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	_, err := io.Copy(pw, pw.resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	fmt.Println("WRITE WAS CALLED YAY!")
	pw.downloaded += len(p)
	fmt.Println(pw.downloaded)
	if pw.total > 0 {
		fmt.Println("Total was greater than 0")
		pw.onProgress(float64(pw.downloaded)/float64(pw.total), percentChannel)
	}
	return len(p), nil
}

func (pw *progressWriter) onProgress(prog float64, percentChan chan float64) {
	progPercent := math.Floor(prog * 100)
	percentChan <- progPercent
}

func WaitForPercentage(percentChan chan float64) tea.Cmd {
	return func() tea.Msg {
		fmt.Println("Wait for percentage was called")
		percentage := <-percentChan
		if percentage == 100 {
			return doneMsg(true)
		}
		return percentageMsg(<-percentChan)
	}
}

func (m Model) Start() tea.Cmd {
	go m.pw.Start()
	return WaitForPercentage(percentChannel)
}

type Model struct {
	bar        progress.Model
	pw         *progressWriter
	url        string
	version    string
	percentage float64
}

func New(version string) Model {
	fmt.Println(version)
	url := fmt.Sprintf("https://go.dev/dl/go%v.linux-amd64.tar.gz", version)
	bar := progress.New()
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
	fmt.Println(resp.ContentLength)
	pw := &progressWriter{
		total: int(resp.ContentLength),
		resp:  resp,
	}
	return Model{
		bar:     bar,
		pw:      pw,
		url:     url,
		version: version,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case percentageMsg:
		m.percentage = float64(msg)
	case doneMsg:
		err := extract.Gz(context.TODO(), m.pw.resp.Body, viper.GetString("installDir"), renamer(m.version))
		if err != nil {
			log.Fatal(err)
		}
		m.pw.resp.Body.Close()
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.bar.ViewAs(m.percentage)
}

func renamer(ver string) extract.Renamer {
	return func(name string) string {
		return strings.Replace(name, "go", ver, 1)
	}
}
