package config

import (
	"github.com/30x/apid"
	"github.com/spf13/viper"
	"log"
	"strings"
)

const (
	configFileName = "apid_config"
)

var cfg *viper.Viper

func GetConfig() apid.ConfigService {
	if cfg == nil {

		viper.SetDefault(configFileName, configFileName)
		configFile := viper.GetString(configFileName)
		configFile = strings.TrimSuffix(configFile, ".yaml")

		cfg = viper.New()

		cfg.SetConfigType("yaml")
		cfg.SetConfigName(configFile)
		cfg.AddConfigPath(".")

		err := cfg.ReadInConfig()
		if err != nil {
			log.Printf("Error in config file '%s': %s", configFileName, err)
		}

		cfg.SetEnvPrefix("apid") // eg. env var "APID_SOMETHING" will bind to config var "something"
		cfg.AutomaticEnv()
	}
	return cfg
}
