package config

import (
	"os"
	"strconv"
)

type Config struct {
	App   AppConfig
	MySQL MySQLConfig
	Redis RedisConfig
	Kafka KafkaConfig
}

type AppConfig struct {
	Mode string
	Port string
}

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers         string
	TopicUserEvents string
	TopicSystemLogs string
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Mode: getEnv("APP_MODE", "debug"),
			Port: getEnv("APP_PORT", "8080"),
		},
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_PASSWORD", "20021008wyj"),
			Database: getEnv("MYSQL_DATABASE", "xx"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Kafka: KafkaConfig{
			Brokers:         getEnv("KAFKA_BROKERS", "localhost:9092"),
			TopicUserEvents: getEnv("KAFKA_TOPIC_USER_EVENTS", "user_events"),
			TopicSystemLogs: getEnv("KAFKA_TOPIC_SYSTEM_LOGS", "system_logs"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
