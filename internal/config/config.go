package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env  string     `yaml:"env" env-default:"local"`
	HTTP HTTPConfig `yaml:"http" env-required:"true"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"1h"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config file is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file " + path + "does not exists")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
