package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
	SSLCert  string `yaml:"sslcert"`
}

func (d Database) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s",
		CONF.Database.Host,
		CONF.Database.Port,
		CONF.Database.User,
		CONF.Database.Password,
		CONF.Database.Database,
		CONF.Database.SSLMode,
		CONF.Database.SSLCert,
	)
}

type Conf struct {
	Database Database `yaml:"db"`
}

func LoadConf() error {
	f, err := os.Open(PATH)
	if err != nil {
		return fmt.Errorf("error opening config: %w", err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&CONF); err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	return nil
}
