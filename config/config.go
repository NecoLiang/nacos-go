// config/config.go

package config

import (
	"github.com/spf13/viper"
)

type NacosConfig struct {
	ServerIP   string ` mapstructure:"server_ip"`
	ServerPort int    `mapstructure:"server_port"`
	Namespace  string `mapstructure:"namespace"`
}

type ServiceConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type AppConfig struct {
	//Nacos   NacosConfig   `mapstructure:"nacos"`
	//Service ServiceConfig `mapstructure:"service"`
	Go struct {
		Nacos   NacosConfig   `yaml:"nacos"`
		Service ServiceConfig `yaml:"service"`
	} `yaml:"go"`
}

func LoadConfig() (AppConfig, error) {
	var config AppConfig
	viper.SetConfigName("application")
	viper.AddConfigPath("resources")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
