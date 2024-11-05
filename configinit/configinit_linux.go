//go:build linux

package configinit

import (
	"os"

	"github.com/spf13/viper"
)

func Init() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configdir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	viper.SetConfigName("gvm")
	viper.SetConfigType("json")
	viper.AddConfigPath(configdir + "/gvm/")
	viper.SetDefault("installDir", homedir+"/.local/share/gvm/versions/")
	err = viper.ReadInConfig()
	return err
}
