package cmd

import (
	"bytes"
	"fmt"
	"io"
	"leomick/gvm/components/downloader"
	"leomick/gvm/tools"
	"leomick/gvm/tools/targz"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type installModel struct {
	downloader downloader.Model
	version    string
	extracting bool
	extracted  bool
	spinner    spinner.Model
}

type extractedMsg bool

var checkmark string = lipgloss.NewStyle().
	Foreground(lipgloss.ANSIColor(10)).
	SetString("✓").String()

func installInitialModel(version string) installModel {
	spin := spinner.New()
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("123"))
	spin.Spinner = spinner.MiniDot
	downloader := downloader.New(tools.GetUrl(version))
	return installModel{
		downloader: downloader,
		version:    version,
		extracting: false,
		extracted:  false,
		spinner:    spin,
	}
}

func (m installModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tea.Batch(m.downloader.Start(), m.spinner.Tick)
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg.(type) {
	case downloader.DoneMsg:
		defer m.downloader.Pw.Resp.Body.Close()
		reader := bytes.NewReader(m.downloader.Pw.Content)
		m.extracting = true
		cmds = append(cmds, extract(m.version, reader))
	case extractedMsg:
		m.extracting = false
		m.extracted = true
		cmds = append(cmds, tea.Quit)
	}
	m.downloader, cmd = m.downloader.Update(msg)
	cmds = append(cmds, cmd)
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m installModel) View() string {
	var s string
	if m.extracting {
		s = fmt.Sprintf("%v Extracting...", m.spinner.View())
	} else if m.extracted {
		s = fmt.Sprintf("%v Successfully extracted!\n", checkmark)
	} else {
		s = m.downloader.View()
	}
	return s
}

func extract(version string, reader io.Reader) tea.Cmd {
	return func() tea.Msg {
		err := targz.ExtractTarGZ(reader, viper.GetString("installDir"), Renamer(version))
		if err != nil {
			log.Fatal(err)
		}
		return extractedMsg(true)
	}
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a specified go version",
	Long: `Installs a specified go version. For example:
"gvm install latest" installs the latest version
"gvm install 1.23.2" installs go version 1.23.2`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ver := args[0]
		if ver == "latest" {
			tbver, err := tools.GetLatestVer()
			if err != nil {
				log.Fatal(err)
			}
			ver = tbver
		}
		_, err := os.Stat(viper.GetString("installDir") + ver)
		switch {
		case os.IsNotExist(err):
			p := tea.NewProgram(installInitialModel(ver))
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		case err != nil:
			log.Fatal(err)
		default:
			fmt.Println("That version is already installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Renamer(ver string) func(string) string {
	return func(name string) string {
		return strings.Replace(name, "go", ver, 1)
	}
}
