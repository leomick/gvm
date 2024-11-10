package cmd

import (
	"fmt"
	"leomick/gvm/components/downloader"
	"leomick/gvm/tools"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type installModel struct {
	downloader downloader.Model
}

func installInitialModel(version string) installModel {
	downloader := downloader.New(version)
	return installModel{
		downloader: downloader,
	}
}

func (m installModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return m.downloader.Start()
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.downloader, cmd = m.downloader.Update(msg)
	return m, cmd
}

func (m installModel) View() string {
	return m.downloader.View()
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
