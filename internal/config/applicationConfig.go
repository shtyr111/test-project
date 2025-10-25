package config

import (
	"os"

	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}

var fullConfig Config
var DATABASE_CONFIG = fullConfig.Database
var SERVER_CONFIG = fullConfig.Server

func LoadApplicationConfig() {
	buf, err := os.ReadFile("./config/application.yml")
	if err != nil {
		log.Errorln("Ошибка при загрузке конфигурации", err.Error())
	}

	var cfg Config
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		log.Errorln("Ошибка при загрузке конфигурации", err.Error())
	}

	fullConfig = cfg
	DATABASE_CONFIG = cfg.Database
	SERVER_CONFIG = cfg.Server
}
