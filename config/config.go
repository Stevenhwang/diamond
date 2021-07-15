package config

import (
	"log"

	"github.com/spf13/viper"
)

var Config *viper.Viper

func init() {
	Config = viper.New()
	Config.SetConfigFile("./config.json")
	err := Config.ReadInConfig()
	if err != nil {
		log.Fatalf("read config failed: %v", err)
	}
}
