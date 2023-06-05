package config

import "os"

const (
	defaultStorageHost     = "127.0.0.1"
	defaultStoragePort     = "5432"
	defaultStorageUser     = "postgres"
	defaultStoragePassword = "password"
	defaultPostDBName      = "posts"
)
const (
	envNameStorageHost     = "STORAGE_HOST"
	envNameStoragePort     = "STORAGE_PORT"
	envNameStorageUser     = "STORAGE_USER"
	envNameStoragePassword = "STORAGE_PASSWORD"
	envNamePostDBName      = "POST_DB_NAME"
)

type Storage struct {
	Host       string
	Port       string
	User       string
	Password   string
	PostDBName string
}

func newDefaultStorageConfig() Storage {
	return Storage{
		Host:       defaultStorageHost,
		Port:       defaultStoragePort,
		User:       defaultStorageUser,
		Password:   defaultStoragePassword,
		PostDBName: defaultPostDBName,
	}
}

func (c *Storage) parseEnv() {
	envHost := os.Getenv(envNameStorageHost)
	if envHost != "" {
		c.Host = envHost
	}

	envPort := os.Getenv(envNameStoragePort)
	if envPort != "" {
		c.Port = envPort
	}

	envUser := os.Getenv(envNameStorageUser)
	if envUser != "" {
		c.User = envUser
	}

	envPassword := os.Getenv(envNameStoragePassword)
	if envPassword != "" {
		c.Password = envPassword
	}

	envPostDBName := os.Getenv(envNamePostDBName)
	if envPostDBName != "" {
		c.PostDBName = envPostDBName
	}
}
