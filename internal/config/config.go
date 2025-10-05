// Package config provides app's config
package config

import (
	"fmt"
	"os"
	"strconv"
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
	JwtSecretKey         []byte
	JwtSigningMethod     *jwt.SigningMethodHMAC
	AccessTokenLifeTime  time.Duration
	RefreshTokenLifeTime time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type Config struct {
	DBConfig
	APIConfig
	TokenConfig
	RedisConfig
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

	accessExp, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	refreshExp, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP"))

	c.TokenConfig = TokenConfig{
		AppName:              os.Getenv("APP_NAME"),
		JwtSecretKey:         []byte(os.Getenv("JWT_SECRET")),
		JwtSigningMethod:     jwt.SigningMethodHS256,
		AccessTokenLifeTime:  time.Duration(accessExp) * time.Minute,
		RefreshTokenLifeTime: time.Duration(refreshExp) * time.Minute,
	}

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	c.RedisConfig = RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	}

	if c.DBConfig.Host == "" ||
		c.DBConfig.Port == "" ||
		c.DBConfig.Username == "" ||
		c.DBConfig.Password == "" ||
		c.APIConfig.APIHost == "" ||
		c.APIConfig.APIPort == "" {
		return fmt.Errorf("there's an empty payload")
	}

	return nil
}

func (c *Config) GetDomain() string {
	return c.APIConfig.Domain
}

func (c *Config) IsSecure() bool {
	return c.APIConfig.IsSecure
}
