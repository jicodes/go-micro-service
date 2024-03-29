package application

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddr string
	ServerPort uint16
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddr: "localhost:6379",
		ServerPort: 8080,
	}

	if redisAddr, exists := os.LookupEnv("REDIS_ADDRESS"); exists {
		cfg.RedisAddr = redisAddr
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		port, err := strconv.ParseUint(serverPort, 10, 16)
		if err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	return cfg
}