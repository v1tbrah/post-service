package config

import "os"

const (
	defaultStorageHost             = "127.0.0.1"
	defaultStoragePort             = "5432"
	defaultStorageUser             = "postgres"
	defaultStoragePassword         = "postgres"
	defaultPostDBName              = "postgres"
	defaultPostTableName           = "post"
	defaultHashtagTableName        = "hashtag"
	defaultHashtagPerPostTableName = "hashtag_per_post"
)
const (
	envNameStorageHost             = "STORAGE_HOST"
	envNameStoragePort             = "STORAGE_PORT"
	envNameStorageUser             = "STORAGE_USER"
	envNameStoragePassword         = "STORAGE_PASSWORD"
	envNamePostDBName              = "POST_DB_NAME"
	envNamePostTableName           = "POST_TABLE_NAME"
	envNameHashtagTableName        = "HASHTAG_TABLE_NAME"
	envNameHashtagPerPostTableName = "HASHTAG_PER_POST_TABLE_NAME"
)

type StorageConfig struct {
	Host                    string
	Port                    string
	User                    string
	Password                string
	PostDBName              string
	PostTableName           string
	HashtagTableName        string
	HashtagPerPostTableName string
}

func newDefaultStorageConfig() StorageConfig {
	return StorageConfig{
		Host:                    defaultStorageHost,
		Port:                    defaultStoragePort,
		User:                    defaultStorageUser,
		Password:                defaultStoragePassword,
		PostDBName:              defaultPostDBName,
		PostTableName:           defaultPostTableName,
		HashtagTableName:        defaultHashtagTableName,
		HashtagPerPostTableName: defaultHashtagPerPostTableName,
	}
}

func (c *StorageConfig) parseEnv() {
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

	envPostTableName := os.Getenv(envNamePostTableName)
	if envPostTableName != "" {
		c.PostTableName = envPostTableName
	}

	envHashtagTableName := os.Getenv(envNameHashtagTableName)
	if envHashtagTableName != "" {
		c.HashtagTableName = envHashtagTableName
	}

	envHashtagPerPostTableName := os.Getenv(envNameHashtagPerPostTableName)
	if envHashtagPerPostTableName != "" {
		c.HashtagPerPostTableName = envHashtagPerPostTableName
	}
}
