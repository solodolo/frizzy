package config

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	DefaultContentDir  string = "content"
	DefaultPagesDir    string = "pages"
	DefaultTemplateDir string = "templates"
)

// Config holds the configuration options for the
// frizzy project
type Config struct {
	RootPath    string
	ContentDir  string
	PagesDir    string
	OutputPath  string
	TemplateDir string
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

func (receiver *Config) GetContentPath() string {
	return filepath.Join(receiver.RootPath, receiver.ContentDir)
}

func (receiver *Config) GetPagesPath() string {
	return filepath.Join(receiver.RootPath, receiver.PagesDir)
}

func (receiver *Config) GetTemplatePath() string {
	return filepath.Join(receiver.RootPath, receiver.TemplateDir)
}

func loadConfigObject(configStream io.Reader) (*Config, error) {
	dec := json.NewDecoder(configStream)

	var c *Config = &Config{}
	if err := dec.Decode(c); err != nil && err != io.EOF {
		return nil, err
	} else if c == nil {
		return nil, errors.New("json config is null")
	}

	if c.ContentDir == "" {
		c.ContentDir = DefaultContentDir
	}

	if c.PagesDir == "" {
		c.PagesDir = DefaultPagesDir
	}

	if c.TemplateDir == "" {
		c.TemplateDir = DefaultTemplateDir
	}

	return c, nil
}
