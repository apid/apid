package config

import (
	"github.com/30x/apid"
	"github.com/spf13/viper"
	"log"
	"strings"
	"os"
)

const (
	localStoragePathKey     = "local_storage_path"
	localStoragePathDefault = "/var/tmp/apid"

	configFileEnvVar = "APID_CONFIG_FILE"

	configFileType    = "yaml"
	configFileNameKey = "apid_config_filename"
	configPathKey     = "apid_config_path"

	defaultConfigFilename = "apid_config.yaml"
	defaultConfigPath     = "."
)

var cfg *viper.Viper

func GetConfig() apid.ConfigService {
	if cfg == nil {

		cfg = viper.New()

		// for config file search path
		cfg.SetConfigType(configFileType)

		cfg.SetDefault(configPathKey, defaultConfigPath)
		configFilePath := cfg.GetString(configPathKey)
		cfg.AddConfigPath(configFilePath)

		cfg.SetDefault(configFileNameKey, defaultConfigFilename)
		configFileName := cfg.GetString(configFileNameKey)
		configFileName = strings.TrimSuffix(configFileName, ".yaml")
		cfg.SetConfigName(configFileName)

		// for user-specified absolute config file
		configFile, ok := os.LookupEnv(configFileEnvVar); if ok {
			cfg.SetConfigFile(configFile)
		}

		cfg.SetDefault(localStoragePathKey, localStoragePathDefault)

		err := cfg.ReadInConfig()
		if err != nil {
			log.Printf("Error in config file '%s': %s", configFileNameKey, err)
		}

		cfg.SetEnvPrefix("apid") // eg. env var "APID_SOMETHING" will bind to config var "something"
		cfg.AutomaticEnv()
	}
	return cfg
}
