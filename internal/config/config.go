package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"` // env-required:"true"
	Storage    `yaml:"storage" env-required:"true"`
	HTTPServer `yaml:"http_server" env-required:"true"`
	Stan       `yaml:"stan" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	DB_Name  string `yaml:"db" env-default:"postgres"`
}

type Stan struct {
	URL       string `yaml:"url" env-default:"http://127.0.0.1:4222"`
	ClusterID string `yaml:"cluster_id" env-default:"L0_cluster"`
	ClientID  string `yaml:"client_id" env-default:"L0_sub"`
	UserCreds string `yaml:"user_creds" env-default:""`
	Channel   string `yaml:"channel" env-default:"L0_chan"`
}

func MustLoad() *Config {
	// configPath := os.Getenv("CONFIG_PATH")
	// configPath := "config/local.yaml"
	configPath := "config/local-docker.yaml"
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
