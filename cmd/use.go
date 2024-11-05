package cmd

import (
	"errors"
	"leomick/gvm/downloader"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Sets the current go version to the specified version",
	Long: `Makes the go command be a specified version. For example:
Running "gvm use 1.23.2" then running "go version" would print "go version go1.23.2 youros/yourcpuarchitecture"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ver := args[0]
		if ver == "latest" {
			tbver, err := downloader.GetLatestVer()
			if err != nil {
				log.Fatal(err)
			}
			ver = tbver
		}
		install, err := cmd.Flags().GetBool("install")
		if err != nil {
			log.Fatal(err)
		}
		_, err = os.Stat(viper.GetString("installDir") + ver)
		if os.IsNotExist(err) {
			if !install {
				log.Fatal(errors.New("You are trying to use a go version that is not installed through gvm"))
			}
			err = downloader.Download(ver)
			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Fatal(err)
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
