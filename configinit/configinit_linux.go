//go:build linux

package configinit

import "github.com/spf13/viper"

func Init() error {
	viper.SetConfigName("gvm")
	viper.SetConfigType("json")
	viper.AddConfigPath("~/.config/gvm/")
	viper.SetDefault("installDir", "~/.local/share/gvm/versions/")
	err := viper.ReadInConfig()
	return err
}
