package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	defaultLogLvl = zerolog.InfoLevel
	envNameLogLvl = "LOG_LVL"
)

type Config struct {
	GRPCConfig    GRPCConfig
	StorageConfig StorageConfig
	KafkaConfig   KafkaConfig

	LogLvl zerolog.Level
}

func NewDefaultConfig() Config {
	return Config{
		GRPCConfig:    newDefaultGRPCConfig(),
		StorageConfig: newDefaultStorageConfig(),
		KafkaConfig:   newDefaultKafkaConfig(),
		LogLvl:        defaultLogLvl,
	}
}

func (c *Config) ParseEnv() error {
	c.GRPCConfig.parseEnv()

	c.StorageConfig.parseEnv()

	c.KafkaConfig.parseEnv()

	if err := c.parseEnvLogLvl(); err != nil {
		return errors.Wrapf(err, "parse log lvl")
	}

	return nil
}

func (c *Config) parseEnvLogLvl() error {
	envLogLvl := os.Getenv(envNameLogLvl)
	if envLogLvl != "" {
		logLevel, err := zerolog.ParseLevel(envLogLvl)
		if err != nil {
			return fmt.Errorf("parse log lvl: %s", envLogLvl)
		}
		c.LogLvl = logLevel
	}
	return nil
}
