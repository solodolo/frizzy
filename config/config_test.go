package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestConfigLoadsCorrectValues(t *testing.T) {
	expectedRoot := "/path/to/project/root"
	expectedContent := "path/to/content"

	configJSON := fmt.Sprintf(`{"RootPath": %q, "ContentPath": %q}`, expectedRoot, expectedContent)
	if config, err := loadConfigObject(strings.NewReader(configJSON)); err != nil {
		t.Errorf("expected no error, got %q", err)
	} else if config.RootPath != expectedRoot {
		t.Errorf(`expected RootPath to be %q, got %q`, expectedRoot, config.RootPath)
	} else if config.ContentPath != expectedContent {
		t.Errorf(`expected ContentPath to be %q got %q`, expectedContent, config.ContentPath)
	}
}
