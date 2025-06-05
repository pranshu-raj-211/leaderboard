package config

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis struct {
		Address    string `yaml:"addr"`
		Password   string `yaml:"password"`
		DB         int    `yaml:"db"`
		MaxRetries int    `yaml:"max_retries"`
	} `yaml:"redis"`

	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`

	Leaderboard struct {
		UpdateIntervalSecs int `yaml:"update_interval_secs"`
		TopPlayersLimit    int `yaml:"top_players_limit"`
		CacheExpiryMins    int `yaml:"cache_expiry_mins"`
	} `yaml:"leaderboard"`
}

var AppConfig *Config
var logger *zap.Logger

func LoadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		Error("Failed to open config file", map[string]any{"err": err})
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		Error("Failed to open config file", map[string]any{"err": err})
		return err
	}

	AppConfig = &cfg
	return nil
}

func Info(msg string, data map[string]any) {
	fields := make([]zap.Field, 0, len(data))
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}
	logger.Info(msg, fields...)
}

func Error(msg string, data map[string]any) error {
	fields := make([]zap.Field, 0, len(data))
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}
	logger.Error(msg, fields...)
	return errors.New(msg)
}

func Fatal(msg string, data map[string]any) {
	fields := make([]zap.Field, 0, len(data))
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}
	logger.Fatal(msg, fields...)
}

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// TODO: move to config file
	logFile := "app.log"
	config.OutputPaths = []string{
		"stdout",
		logFile,
	}
	config.ErrorOutputPaths = []string{
		"stderr",
		logFile,
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	return logger
}
