// Package config provides app's config
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Driver   string
}

type APIConfig struct {
	APIHost  string
	APIPort  string
	Domain   string
	IsSecure bool
}

type TokenConfig struct {
	AppName              string
	JwtSignaturKey       []byte
	JwtSigningMethod     *jwt.SigningMethodHMAC
	AccessTokenLifeTime  time.Duration
	RefreshTokenLifeTime time.Duration
}

type Config struct {
	DBConfig
	APIConfig
	TokenConfig
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cfg.readConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil

}

func (c *Config) readConfig() error {
	c.DBConfig = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Driver:   os.Getenv("DB_DRIVER"),
	}
	isSecure := os.Getenv("ENV") == "production"

	c.APIConfig = APIConfig{
		APIHost:  os.Getenv("API_HOST"),
		APIPort:  os.Getenv("API_PORT"),
		Domain:   os.Getenv("DOMAIN"),
		IsSecure: isSecure,
	}

	accessTokenLifeTime := time.Duration(1) * time.Hour
	refreshTokenLifeTime := time.Duration(24*7) * time.Hour

	c.TokenConfig = TokenConfig{
		AppName:              os.Getenv("APP_NAME"),
		JwtSignaturKey:       []byte(os.Getenv("TOKEN_KEY")),
		JwtSigningMethod:     jwt.SigningMethodHS256,
		AccessTokenLifeTime:  accessTokenLifeTime,
		RefreshTokenLifeTime: refreshTokenLifeTime,
	}

	if c.Host == "" ||
		c.Port == "" ||
		c.Username == "" ||
		c.Password == "" ||
		c.APIHost == "" ||
		c.APIPort == "" {
		return fmt.Errorf("there's an empty payload")
	}

	return nil
}

func (c *Config) GetDomain() string {
	return c.Domain
}

func (c *Config) IsSecure() bool {
	return c.APIConfig.IsSecure
}
