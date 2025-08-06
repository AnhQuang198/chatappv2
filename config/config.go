package config

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type Config struct {
	//load config server application from file yaml
	Server struct {
		Port string
	}

	//load config for database from file yaml
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		Driver   string
		SSLMode  string
	}
}

func LoadApplicationConfig() (*Config, error) {
	filePath := fmt.Sprintf("config/resource/application.yaml")
	content, err := os.ReadFile(filePath) //load all properties from yaml file
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}

	expanded := os.ExpandEnv(string(content)) //replace value from environment

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(strings.NewReader(expanded)); err != nil {
		return nil, fmt.Errorf("viper read config error: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil { //mapping config from yaml to Config struct
		return nil, fmt.Errorf("viper unmarshal error: %w", err)
	}

	return &cfg, nil
}

// InitPostgresConnection init connection to postgres with config from yaml file
func InitPostgresConnection(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	return db, nil
}
