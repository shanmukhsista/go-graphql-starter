package config

import (
	"flag"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

// MustGetString Get a String Value from config.
func MustGetString(key string) string {
	configValue := viper.GetString(key)
	if len(strings.TrimSpace(configValue)) == 0 {
		panic(fmt.Sprintf("unable to load config key %s", key))
	}
	return configValue
}

// MustGetStringSet Get a unique set of Values from configuration. Dedupe input list for
// uniqueness.
func MustGetStringSet(key string) []string {
	valuesList := viper.GetStringSlice(key)
	if len(valuesList) == 0 {
		panic(fmt.Sprintf("unable to load config key %s", key))
	}
	return lo.Uniq(valuesList)
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (err error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func MustGetConfigPathFromFlags(configKey string) string {
	var configPath *string
	configPath = flag.String(configKey, "", "Config Path")
	flag.Parse()
	log.Debug().Msgf("Using config path %s", &configPath)
	if configPath == nil || lo.IsEmpty(*configPath) {
		log.Fatal().Msgf("Unable to load config path. Empty Path specified. ")
	}
	return *configPath
}

func MustLoadConfigAtPath(configPath string) error {
	err := LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	return err
}
