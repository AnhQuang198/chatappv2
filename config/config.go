package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	cfg     *Config
	cfgOnce sync.Once
	cfgErr  error
)

type Config struct {
	//load config server application from file yaml
	Server struct {
		Port        string
		AllowOrigin string
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
		MaxIdle  int
		MaxConn  int
		MaxLife  int
	}
}

func LoadApplicationConfig() (*Config, error) {
	cfgOnce.Do(func() {
		filePath := fmt.Sprintf("config/resource/application.yaml")
		content, err := os.ReadFile(filePath) //load all properties from yaml file
		if err != nil {
			cfgErr = fmt.Errorf("read config file error: %w", err)
			return
		}

		expanded := os.ExpandEnv(string(content)) //replace value from environment

		viper.SetConfigType("yaml")
		if err := viper.ReadConfig(strings.NewReader(expanded)); err != nil {
			cfgErr = fmt.Errorf("viper read config error: %w", err)
			return
		}

		var result Config
		if err := viper.Unmarshal(&result); err != nil { //mapping config from yaml to Config struct
			cfgErr = fmt.Errorf("viper unmarshal config error: %w", err)
			return
		}
		cfg = &result
	})

	return cfg, cfgErr
}

// InitPostgresConnection init connection to postgres with config from yaml file
func InitPostgresConnection() (*sql.DB, error) {
	cfg, err := LoadApplicationConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
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
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}

	//config connection pool
	db.SetMaxOpenConns(cfg.Database.MaxConn)
	db.SetMaxIdleConns(cfg.Database.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLife) * time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping db: %w", err)
	}
	return db, nil
}
