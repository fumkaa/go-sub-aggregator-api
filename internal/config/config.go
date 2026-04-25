package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string           `yaml:"env" env-default:"local"`
	Storage    StorageConfig    `yaml:"storage"`
	HttpServer HttpServerConfig `yaml:"http_server"`
}

type HttpServerConfig struct {
	Port         int           `yaml:"port" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"30s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"30s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type StorageConfig struct {
	Driver            string        `yaml:"driver"`
	Host              string        `yaml:"host" env-default:"localhost"`
	Port              int           `yaml:"port" env-default:"5432"`
	DBName            string        `yaml:"db_name"`
	User              string        `env:"DB_USER" env-required:"true"`
	Password          string        `env:"DB_PASSWORD" env-required:"true"`
	MaxConns          int32         `yaml:"max_conns" env-default:"25"`
	MinConns          int32         `yaml:"min_conns" env-default:"5"`
	MaxConnIdleTime   time.Duration `yaml:"max_conn_idle_time" env-default:"30m"`
	HealthCheckPeriod time.Duration `yaml:"health_check_period" env-default:"1m"`
}

func (s *StorageConfig) DSN() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s", s.Driver, s.User, s.Password, s.Host, s.Port, s.DBName)
}

func MustLoad() *Config {
	if err := godotenv.Load(".env"); err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is required")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read env: " + err.Error())
	}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config file: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "config-path", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
