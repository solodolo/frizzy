package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

var expectedRoot = "/path/to/project/root"

func TestConfigSetsDefaultContentDir(t *testing.T) {
	expectedContent := filepath.Join(expectedRoot, "content")

	configJSON := fmt.Sprintf(`{"RootPath": %q}`, expectedRoot)
	if config, err := loadConfigObject(strings.NewReader(configJSON)); err != nil {
		t.Errorf("expected no error, got %q", err)
	} else if config.RootPath != expectedRoot {
		t.Errorf(`expected RootPath to be %q, got %q`, expectedRoot, config.RootPath)
	} else if config.GetContentPath() != expectedContent {
		t.Errorf(`expected content path to be %q got %q`, expectedContent, config.GetContentPath())
	}
}

func TestConfigSetsDefaultPageDir(t *testing.T) {
	expectedPages := filepath.Join(expectedRoot, "pages")

	configJSON := fmt.Sprintf(`{"RootPath": %q}`, expectedRoot)
	if config, err := loadConfigObject(strings.NewReader(configJSON)); err != nil {
		t.Errorf("expected no error, got %q", err)
	} else if config.RootPath != expectedRoot {
		t.Errorf(`expected RootPath to be %q, got %q`, expectedRoot, config.RootPath)
	} else if config.GetPagesPath() != expectedPages {
		t.Errorf(`expected pages path to be %q got %q`, expectedPages, config.GetContentPath())
	}
}

func TestConfigSetsCorrectContentDir(t *testing.T) {
	dir := "/some/content/dir"
	expected := filepath.Join(expectedRoot, dir)

	configJSON := fmt.Sprintf(`{"RootPath": %q, "ContentDir": %q}`, expectedRoot, dir)
	if config, err := loadConfigObject(strings.NewReader(configJSON)); err != nil {
		t.Errorf("expected no error, got %q", err)
	} else if config.RootPath != expectedRoot {
		t.Errorf(`expected RootPath to be %q, got %q`, expectedRoot, config.RootPath)
	} else if config.GetContentPath() != expected {
		t.Errorf(`expected content path to be %q got %q`, expected, config.GetContentPath())
	}
}

func TestConfigSetsCorrectPagesDir(t *testing.T) {
	dir := "/some/pages/dir"
	expected := filepath.Join(expectedRoot, dir)

	configJSON := fmt.Sprintf(`{"RootPath": %q, "PagesDir": %q}`, expectedRoot, dir)
	if config, err := loadConfigObject(strings.NewReader(configJSON)); err != nil {
		t.Errorf("expected no error, got %q", err)
	} else if config.RootPath != expectedRoot {
		t.Errorf(`expected RootPath to be %q, got %q`, expectedRoot, config.RootPath)
	} else if config.GetPagesPath() != expected {
		t.Errorf(`expected content path to be %q got %q`, expected, config.GetContentPath())
	}
}
