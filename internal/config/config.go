package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf `json:"logger"`
	DB     DBConf     `json:"db"`
	HTTP   HTTPConf   `json:"http"`
	AMPQ   AMPQConf   `json:"ampq"`
}

type LoggerConf struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type DBConf struct {
	ConnectionString string `json:"connection_string"`
}

type HTTPConf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	GrpcPort string `json:"grpc_port"`
}

type AMPQConf struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

func New(configFile string) (Config, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	return Config{
		LoggerConf{Level: viper.GetString("logger.level"), File: viper.GetString("logger.file")},
		DBConf{ConnectionString: viper.GetString("db.connection_string")},
		HTTPConf{Host: viper.GetString("http.host"), Port: viper.GetString("http.port"), GrpcPort: viper.GetString("http.grpc_port")},
		AMPQConf{URI: viper.GetString("ampq.uri"), Name: viper.GetString("ampq.name")},
	}, nil
}
