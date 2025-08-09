package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	DataBase struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	RedisConfig struct {
		Host     string        `yaml:"host"`
		Port     int           `yaml:"port"`
		Password string        `yaml:"password"`
		DB       int           `yaml:"db"`
		TTL      time.Duration `yaml:"ttl"`
	} `yaml:"redis"`

	KafkaConfig struct {
		Broker  string `yaml:"broker"`
		Topic   string `yaml:"topic"`
		GroupID string `yaml:"group_id"`
	} `yaml:"kafka"`

	LoggerConfig struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logger"`
}

func MustLoad() *Config {
	path := FetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file is not exist :" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func FetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
