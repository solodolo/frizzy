package file

import (
	"os"
	"path/filepath"

	"mettlach.codes/frizzy/config"
)

func GetContentPaths(subpath string) []string {
	contentDirs := []string{}
	// assumes GetContentPath returns an absolute path
	walkPath := filepath.Join(config.GetContentPath(), subpath)

	filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		contentDirs = append(contentDirs, path)
		return nil
	})

	return contentDirs
}