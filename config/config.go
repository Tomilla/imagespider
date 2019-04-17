package config

import "os"

type Config struct {
	Path           string
	SeperateFolder bool
}

var DefaultConfig = getDefaultConfig()

func getDefaultConfig() *Config {
	if wd, err := os.Getwd(); err == nil {
		return &Config{wd + "fruit", false}
	} else {
		panic(err)
	}
}
