package file

import (
	"os"
	"path/filepath"
	"strings"

	"mettlach.codes/frizzy/config"
)

// GetPathFunc allows normal functions to be used
// to retreive paths from various sources
type GetPathFunc func(subpath string) []string

// GetContentPaths returns an array of paths to each file
// in the project's <contentroot>/<subpath> directory
func GetContentPaths(subpath string) []string {
	contentDirs := []string{}
	config := config.GetLoadedConfig()
	// assumes GetContentPath returns an absolute path
	walkPath := filepath.Join(config.GetContentPath(), subpath)

	filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		contentDirs = append(contentDirs, path)
		return nil
	})

	return contentDirs
}

func GetRelativePathTo(filepath string) string {
	config := config.GetLoadedConfig()
	return strings.TrimPrefix(filepath, config.RootPath)
}
