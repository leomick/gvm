package tools

import (
	"os"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/viper"
)

func GetVersions() ([]*semver.Version, error) {
	var versions []*semver.Version
	installDir := viper.GetString("installDir")
	_, err := os.Stat(installDir)
	if os.IsNotExist(err) {
		return []*semver.Version{}, nil
	} else if err != nil {
		return []*semver.Version{}, err
	}
	contents, err := os.ReadDir(installDir)
	if err != nil {
		return []*semver.Version{}, err
	}
	for _, c := range contents {
		if c.IsDir() {
			version, err := semver.NewVersion(c.Name())
			if err != nil {
				return []*semver.Version{}, err
			}
			versions = append(versions, version)
		}
	}
	sort.Sort(sort.Reverse(semver.Collection(versions))) // Use sort.Reverse to turn the ascending sort interface into a descending one.
	return versions, nil
}
