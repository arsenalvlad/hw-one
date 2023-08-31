package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Server  ServerConf  `yaml:"server"`
	Storage StorageConf `yaml:"storage"`
}

type StorageConf struct {
	Type   string `yaml:"type"`
	Memory struct {
		DeleteTime time.Duration `yaml:"deleteTime"`
	} `yaml:"memory"`
	Postgres Postgres `yaml:"postgres"`
}

type Postgres struct {
	Host      string `yaml:"host"`
	Port      int64  `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
	SSLMode   string `yaml:"sslmode"`
	Migration struct {
		Path  string `yaml:"path"`
		Table string `yaml:"table"`
	} `yaml:"migration"`
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

func NewConfig(configPath string) Config {
	var conf Config

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return conf
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.Database, p.SSLMode)
}

func (p *Postgres) MigrateDSN() string {
	return fmt.Sprintf("%s&x-migrations-table=%s", p.DSN(), p.Migration.Table)
}

func (s *ServerConf) Address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
