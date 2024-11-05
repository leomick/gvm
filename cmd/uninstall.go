package cmd

import (
	"fmt"
	"leomick/gvm/downloader"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstalls a specified version that has been previously installed with gvm",
	Long:  `Uninstalls a specified version that has been previously installed with gvm`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		installDir := viper.GetString("installDir")
		ver := args[0]
		if ver == "latest" {
			tbver, err := downloader.GetLatestVer()
			if err != nil {
				log.Fatal(err)
			}
			ver = tbver
		}
		_, err := os.Stat(installDir + ver)
		switch {
		case os.IsNotExist(err):
			fmt.Println("You are trying to uninstall a go version that is not installed through gvm")
			os.Exit(1)
		case err != nil:
			log.Fatal(err)
		default:
			err = os.RemoveAll(installDir + ver)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Successfully uninstalled go version %v. Remember to use the \"gvm use\" command to set another version\n", ver)
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
