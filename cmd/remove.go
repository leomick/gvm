package cmd

import (
	"fmt"
	"leomick/gvm/tools"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a specified version that has been previously installed with gvm",
	Long:  `Removes a specified version that has been previously installed with gvm`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		installDir := viper.GetString("installDir")
		ver := args[0]
		if ver == "latest" {
			tbver, err := tools.GetLatestVer()
			if err != nil {
				log.Fatal(err)
			}
			ver = tbver
		}
		_, err := os.Stat(installDir + ver)
		switch {
		case os.IsNotExist(err):
			fmt.Println("You are trying to remove a go version that is not installed through gvm")
			os.Exit(1)
		case err != nil:
			log.Fatal(err)
		default:
			err = os.RemoveAll(installDir + ver)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Successfully removed go version %v. Remember to use the \"gvm use\" command to set another version\n", ver)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
