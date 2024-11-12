package config

import (
	"flag"
	"github.com/spf13/viper"
	"log/slog"
)

type GraphAPIConfig struct {
	HTTPAddr    string
	SecretToken string
	Database    *Database
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func NewGraphAPIConfig() (*GraphAPIConfig, error) {
	v := viper.GetViper()

	configPath := "config/config.yaml"
	path := flag.String("config", configPath, "path to config.yaml config")
	flag.Parse()

	baseConfig, err := newConfig(v, *path)
	if err != nil {
		slog.Error("failed to create config on <NewGraphAPIConfig>", err)
		return nil, ErrConfig
	}

	c := &GraphAPIConfig{
		HTTPAddr:    v.GetString("graph_api.http_addr"),
		SecretToken: v.GetString("graph_api.secret_token"),
		Database:    baseConfig.NewDatabase(),
	}
	slog.Info("Graph API config was created!")
	return c, nil
}

func (c *config) NewDatabase() *Database {
	return &Database{
		Host:     c.Viper.GetString("database.host"),
		Port:     c.Viper.GetInt("database.port"),
		User:     c.Viper.GetString("database.user"),
		Password: c.Viper.GetString("database.password"),
		Name:     c.Viper.GetString("database.db_name"),
	}
}

type config struct {
	Viper *viper.Viper
}

func newConfig(v *viper.Viper, configFile string) (*config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		slog.Error("failed to read a config file on <NewConfig>", err)
		return nil, ErrConfig
	}

	return &config{Viper: v}, nil
}
