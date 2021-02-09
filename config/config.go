package config

import (
	"path/filepath"
)

// GetRootPath returns the full path to the frizzy root
func GetRootPath() string {
	return ""
}

// GetContentPath returns the full path to the content
// dir in the frizzy root
func GetContentPath() string {
	return filepath.Join(GetRootPath(), "content")
}
