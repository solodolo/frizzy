package config

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

// Config holds the configuration options for the
// frizzy project
type Config struct {
	RootPath    string
	ContentPath string
	OutputPath  string
}

var loadedConfig *Config

// LoadConfig reads the json file at configPath and stores
// the read config in loadedConfig
func LoadConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	defer configFile.Close()

	if err != nil {
		return nil, err
	}

	config, loadErr := loadConfigObject(configFile)
	loadedConfig = config

	return loadedConfig, loadErr
}

// GetLoadedConfig returns the cached config in loadedConfig
// Requires that loadedConfig is not nil which can be done
// by calling LoadConfig
func GetLoadedConfig() *Config {
	if loadedConfig == nil {
		log.Fatal("configuration not loaded")
	}

	return loadedConfig
}

func loadConfigObject(configStream io.Reader) (*Config, error) {
	dec := json.NewDecoder(configStream)

	var c *Config = &Config{}
	if err := dec.Decode(c); err != nil && err != io.EOF {
		return nil, err
	} else if c == nil {
		return nil, errors.New("json config is null")
	}

	return c, nil
}
