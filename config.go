package main

import "github.com/golobby/env/v2"

type ServerConfig struct {
	Address string `env:"SERVER_ADDRESS"`
	Port    uint16 `env:"SERVER_PORT"`
}

type CloudflareConfig struct {
	Ttl   uint32 `env:"CF_TTL"`
	KeyId string `env:"CF_TURN_KEY_ID"`
}

type Config struct {
	Server     ServerConfig
	Cloudflare CloudflareConfig
}

func loadConfig() *Config {
	config := &Config{
		Server: ServerConfig{
			Address: "0.0.0.0",
			Port:    1337,
		},
		Cloudflare: CloudflareConfig{
			Ttl: 86400,
		},
	}

	env.Feed(config)

	return config
}
