package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type InvoiceConf struct {
	Title string   `yaml:"title"`
	Font  string   `yaml:"font"`
	Terms []string `yaml:"terms"`
}

type DBConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
	SSLCert  string `yaml:"sslcert"`
}

func (d DBConf) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s",
		CONF.DBConf.Host,
		CONF.DBConf.Port,
		CONF.DBConf.User,
		CONF.DBConf.Password,
		CONF.DBConf.Database,
		CONF.DBConf.SSLMode,
		CONF.DBConf.SSLCert,
	)
}

type Conf struct {
	InvoiceConf `yaml:"invoice"`
	DBConf      DBConf `yaml:"db"`
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
