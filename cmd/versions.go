package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// versionsCmd represents the versions command
var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Lists the go versions installed with gvm",
	Long:  `Lists the go versions installed with gvm`,
	Run: func(cmd *cobra.Command, args []string) {
		var versions []string
		installDir := viper.GetString("installDir")
		_, err := os.Stat(installDir)
		if os.IsNotExist(err) {
			fmt.Println("You have no versions installed")
		}
		contents, err := os.ReadDir(installDir)
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range contents {
			if c.IsDir() {
				versions = append(versions, c.Name())
			}
		}
		if len(contents) == 0 {
			fmt.Println("You have no versions installed")
		} else {
			list := strings.Join(versions, "\n")
			fmt.Printf("You have the following versions installed:\n%v\n", list)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
