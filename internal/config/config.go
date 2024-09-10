package config

import (
	"eth-parser/common"
	"os"
)

type Config struct {
	ServerAddress string
	EthNodeURL    string
}

func Load() *Config {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		EthNodeURL:    getEnv("ETH_NODE_URL", common.CloudFlareRpcUrl),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
