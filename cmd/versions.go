package cmd

import (
	"fmt"
	"leomick/gvm/tools"
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
		versions, err := getVersions()
		if err != nil {
			log.Fatal(err)
		}
		if len(versions) == 0 {
			fmt.Println("You have no versions installed")
		} else {
			sortedVersions, err := tools.SortVersions(versions)
			if err != nil {
				log.Fatal(err)
			}
			list := strings.Join(sortedVersions, "\n")
			fmt.Printf("You have the following versions installed:\n%v\n", list)
		}
	},
}

func getVersions() ([]string, error) {
	var versions []string
	installDir := viper.GetString("installDir")
	_, err := os.Stat(installDir)
	if os.IsNotExist(err) {
		return []string{}, nil
	} else if err != nil {
		return []string{}, err
	}
	contents, err := os.ReadDir(installDir)
	if err != nil {
		return []string{}, err
	}
	for _, c := range contents {
		if c.IsDir() {
			versions = append(versions, c.Name())
		}
	}
	return versions, nil
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
