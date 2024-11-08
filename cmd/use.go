package cmd

import (
	"fmt"
	"leomick/gvm/tools"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var selected = lipgloss.NewStyle().
	BorderLeft(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderLeftForeground(lipgloss.ANSIColor(10))

type useModel struct {
	versions         []string
	filteredversions []string
	cursor           int
	searchBar        textinput.Model
}

func initialModel() useModel {
	textInput := textinput.New()
	textInput.Prompt = ""
	textInput.Placeholder = "Search for a version"
	textInput.CharLimit = 10
	textInput.Width = 20
	versions, err := tools.GetVersions()
	if err != nil {
		log.Fatal(err)
	}
	var stringVersions []string
	for _, v := range versions {
		stringVersions = append(stringVersions, v.Original())
	}
	return useModel{
		versions:  stringVersions,
		cursor:    0,
		searchBar: textInput,
	}
}

func (m useModel) Init() tea.Cmd {
	if len(m.versions) == 0 {
		fmt.Println("You have no versions installed")
		return tea.Quit
	}
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m useModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > -1 {
				m.cursor--
				if m.cursor == -1 {
					cmd = m.searchBar.Focus()
					cmds = append(cmds, cmd)
				}
			}
		case "down":
			if m.cursor < len(m.filteredversions)-1 {
				m.cursor++
				if m.cursor != -1 && m.searchBar.Focused() {
					m.searchBar.Blur()
				}
			}
		case "enter":
			if !m.searchBar.Focused() {
				// make it actually set it later
				return m, tea.Quit
			}
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	m.searchBar, cmd = m.searchBar.Update(msg)
	cmds = append(cmds, cmd)
	m.filteredversions = []string{}
	for _, v := range m.versions {
		if strings.Contains(v, m.searchBar.Value()) {
			m.filteredversions = append(m.filteredversions, v)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m useModel) View() string {
	s := m.searchBar.View() + "\n"
	if len(m.filteredversions) == 0 {
		s += "No results\n"
	} else {
		for i, version := range m.filteredversions {
			if m.cursor == i {
				s += selected.Render(version) + "\n"
			} else {
				s += fmt.Sprintf(" %v\n", version)
			}
		}
	}
	return s
}

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Sets the current go version to the specified version",
	Long: `Makes the go command be a specified version. For example:
Running "gvm use 1.23.2" then running "go version" would print "go version go1.23.2 youros/yourcpuarchitecture"`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		install, err := cmd.Flags().GetBool("install")
		if err != nil {
			log.Fatal(err)
		}
		if len(args) == 1 {
			ver := args[0]
			if ver == "latest" {
				tbver, err := tools.GetLatestVer()
				if err != nil {
					log.Fatal(err)
				}
				ver = tbver
			}
			_, err = os.Stat(viper.GetString("installDir") + ver)
			if os.IsNotExist(err) {
				if !install {
					fmt.Println("You are trying to use a go version that is not installed through gvm")
					os.Exit(1)
				}
				err = tools.Download(ver)
				if err != nil {
					log.Fatal(err)
				}
			} else if err != nil {
				log.Fatal(err)
			}
		} else if install {
			fmt.Println("You need to specify a version when using the install flag")
			os.Exit(1)
		} else {
			p := tea.NewProgram(initialModel())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.
	useCmd.PersistentFlags().Bool("install", false, "installs the specified version if it is not present")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
