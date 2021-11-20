package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfig(env string) (*Config, error) {
	var cfg Config

	file := fmt.Sprintf("config.%s.yaml", env)

	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.SetConfigFile(file)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
