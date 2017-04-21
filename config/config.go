package config

import (
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init loads the environment configuration into
// a viper configuration struct
func Init(env string) {

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(env)
	v.AddConfigPath("config/")

	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	config = v
}

// GetConfig returns the current configuration for the app
// ensure you call init before using this function
func GetConfig() *viper.Viper {
	return config
}
