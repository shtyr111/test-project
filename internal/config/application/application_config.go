package application

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
		Pool     struct {
			MinSize int32 `yaml:"min-size"`
			MaxSize int32 `yaml:"max-size"`
		} `yaml:"pool"`
	} `yaml:"database"`

	Properties struct {
		SendUserToOlbSchedulerCron                 string `yaml:"send-user-to-olb-scheduler-cron"`
		SendUserToOlbSchedulerSectionAdvisoryCron  int    `yaml:"send-user-to-olb-scheduler-section-advisory-lock"`
		SendUserToOlbSchedulerParallelCurrencySend int    `yaml:"send-user-to-olb-scheduler-parallel-currency-send"`
		InternalServerGrpcAddress                  string `yaml:"internal-server-grpc-address"`
	} `yaml:"configuration-properties"`

	Redis struct {
		Address  string `yaml:"address"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		Pool     struct {
			MinSize int `yaml:"min-size"`
			MaxSize int `yaml:"max-size"`
		} `yaml:"pool"`
	} `yaml:"redis"`
}

var fullConfig Config
var DATABASE_CONFIG = fullConfig.Database
var SERVER_CONFIG = fullConfig.Server

func LoadApplicationConfig() Config {
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

	return cfg
}
