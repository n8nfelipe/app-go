package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DB DbConfig `mapstructure:"db"`
}

type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

func keys() []string {
	return []string{
		"db.host",
		"db.port",
		"db.user",
		"db.password",
		"db.dbname",
	}
}

func Load() (*Config, error) {
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5432")
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.dbname", "mydb")

	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, k := range keys() {
		if err := viper.BindEnv(k); err != nil {
			return nil, fmt.Errorf("vinculando env %s: %w", k, err)
		}
	}

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("lendo arquivo .env: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("decodificando config: %w", err)
	}

	return &cfg, nil
}

func GetDbConfig() *DbConfig {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("erro ao carregar config: %v", err))
	}
	return &cfg.DB
}
