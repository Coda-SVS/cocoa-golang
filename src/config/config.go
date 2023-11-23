package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Config = viper.GetViper()

func init() {
	isCreated := false
	if _, err := os.Stat("config.yaml"); errors.Is(err, os.ErrNotExist) {
		os.Create("config.yaml")
		isCreated = true
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	setupConfig()

	if isCreated {
		WriteConfig()
	}
}

func ReadConfig() {
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func setupConfig() {
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
}

func WriteConfig() {
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}
}
