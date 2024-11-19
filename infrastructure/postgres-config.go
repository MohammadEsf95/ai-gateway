package infrastructure

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Port     int    `yaml:"port"`
}

func (conf *PostgresConfig) getConfig() error {
	yml, err := os.ReadFile("infrastructure/config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}

	err = yaml.Unmarshal(yml, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}

	return nil
}

func ConnectPostgres() (*gorm.DB, error) {
	var conf PostgresConfig
	err := conf.getConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
