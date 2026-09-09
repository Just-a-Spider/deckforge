package assets

import (
	"embed"
	"fmt"
)

//go:embed embedded/*
var EmbeddedFiles embed.FS

// LoadCoreAsset retrieves a core engine asset from embedded storage or returns error
func LoadCoreAsset(path string) (string, error) {
	fullPath := fmt.Sprintf("embedded/%s", path)
	content, err := EmbeddedFiles.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("embedded asset not found: %s (%v)", fullPath, err)
	}
	return string(content), nil
}
