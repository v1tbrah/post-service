package config

import (
	"os"
)

const (
	defaultKafkaEnable      = true
	defaultKafkaHost        = "127.0.0.1"
	defaultKafkaPort        = "9092"
	defaultTopicPostCreated = "post_created"
	defaultTopicPostDeleted = "post_deleted"
)
const (
	envNameKafkaEnable      = "KAFKA_ENABLE"
	envNameKafkaHost        = "KAFKA_HOST"
	envNameKafkaPort        = "KAFKA_PORT"
	envNameTopicPostCreated = "TOPIC_POST_CREATED"
	envNameTopicPostDeleted = "TOPIC_POST_DELETED"
)

type Kafka struct {
	Enable bool

	Host string
	Port string

	TopicPostCreated string

	TopicPostDeleted string
}

func newDefaultKafkaConfig() Kafka {
	return Kafka{
		Enable:           defaultKafkaEnable,
		Host:             defaultKafkaHost,
		Port:             defaultKafkaPort,
		TopicPostCreated: defaultTopicPostCreated,
		TopicPostDeleted: defaultTopicPostDeleted,
	}
}

func (c *Kafka) parseEnv() {
	envKafkaEnable := os.Getenv(envNameKafkaEnable)
	if envKafkaEnable != "" {
		c.Enable = envKafkaEnable == "true"
	}

	envKafkaHost := os.Getenv(envNameKafkaHost)
	if envKafkaHost != "" {
		c.Host = envKafkaHost
	}

	envKafkaPort := os.Getenv(envNameKafkaPort)
	if envKafkaPort != "" {
		c.Port = envKafkaPort
	}

	envTopicPostCreated := os.Getenv(envNameTopicPostCreated)
	if envTopicPostCreated != "" {
		c.TopicPostCreated = envTopicPostCreated
	}

	envTopicPostDeleted := os.Getenv(envNameTopicPostDeleted)
	if envTopicPostDeleted != "" {
		c.TopicPostDeleted = envTopicPostDeleted
	}
}
