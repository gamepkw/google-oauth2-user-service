package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitializeViper() {

	viper.SetConfigName("config")
	viper.AddConfigPath("../config")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}
