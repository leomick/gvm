package cmd

import (
	"fmt"
	"leomick/gvm/tools"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// versionsCmd represents the versions command
var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Lists the go versions installed with gvm",
	Long:  `Lists the go versions installed with gvm`,
	Run: func(cmd *cobra.Command, args []string) {
		versions, err := tools.GetVersions()
		if err != nil {
			log.Fatal(err)
		}
		if len(versions) == 0 {
			fmt.Println("You have no versions installed")
		} else {
			var stringVersions []string
			for _, v := range versions {
				stringVersions = append(stringVersions, v.Original())
			}
			if err != nil {
				log.Fatal(err)
			}
			list := strings.Join(stringVersions, "\n")
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
