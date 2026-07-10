/*
Copyright 2026 stonebirdjx.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary yaml config file
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	yamlContent := `
log:
  level: "info"
  format: "console"
`
	tmpFile := filepath.Join(tmpDir, "server.yaml")
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	cfg, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Load returned unexpected error: %v", err)
	}

	if cfg.Log.Level != "info" {
		t.Errorf("expected Log.Level to be 'info', got '%s'", cfg.Log.Level)
	}
	if cfg.Log.Format != "console" {
		t.Errorf("expected Log.Format to be 'console', got '%s'", cfg.Log.Format)
	}
}
