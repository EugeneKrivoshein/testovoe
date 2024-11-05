package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database   DatabaseConfig `json:"database"`
	AdminToken string         `json:"admin_token"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

// Загрузка конфигурации из файла
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Получение административного токена
func GetAdminToken() string {
	cfg, err := LoadConfig("config/config.json")
	if err != nil {
		return ""
	}
	return cfg.AdminToken
}
