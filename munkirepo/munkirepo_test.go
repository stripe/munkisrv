package munkirepo

import (
	"testing"
)

func TestRepoEmbedded(t *testing.T) {
	// Test that the embedded filesystem is accessible
	entries, err := Repo.ReadDir(".")
	if err != nil {
		t.Fatalf("Failed to read embedded repo root: %v", err)
	}

	// Check that we have the expected directories
	expectedDirs := map[string]bool{
		"catalogs":         false,
		"client_resources": false,
		"icons":            false,
		"manifests":        false,
	}

	for _, entry := range entries {
		if entry.IsDir() {
			expectedDirs[entry.Name()] = true
		}
	}

	for dir, found := range expectedDirs {
		if !found {
			t.Errorf("Expected directory '%s' not found in embedded repo", dir)
		}
	}
}

func TestRepoCatalogsAccessible(t *testing.T) {
	// Test that catalogs directory is accessible
	entries, err := Repo.ReadDir("catalogs")
	if err != nil {
		t.Fatalf("Failed to read catalogs directory: %v", err)
	}

	// Should have at least one catalog file
	if len(entries) == 0 {
		t.Error("Expected at least one catalog file")
	}

	// Check for the 'all' catalog specifically
	foundAll := false
	for _, entry := range entries {
		if entry.Name() == "all" && !entry.IsDir() {
			foundAll = true
			break
		}
	}

	if !foundAll {
		t.Error("Expected 'all' catalog file not found")
	}
}

func TestRepoManifestsAccessible(t *testing.T) {
	// Test that manifests directory is accessible
	entries, err := Repo.ReadDir("manifests")
	if err != nil {
		t.Fatalf("Failed to read manifests directory: %v", err)
	}

	// Should have at least one manifest
	if len(entries) == 0 {
		t.Error("Expected at least one manifest file")
	}
}

func TestRepoIconsAccessible(t *testing.T) {
	// Test that icons directory is accessible
	entries, err := Repo.ReadDir("icons")
	if err != nil {
		t.Fatalf("Failed to read icons directory: %v", err)
	}

	// Should have at least one icon file
	if len(entries) == 0 {
		t.Error("Expected at least one icon file")
	}
}

func TestRepoClientResourcesAccessible(t *testing.T) {
	// Test that client_resources directory is accessible
	entries, err := Repo.ReadDir("client_resources")
	if err != nil {
		t.Fatalf("Failed to read client_resources directory: %v", err)
	}

	// Should have at least one client resource
	if len(entries) == 0 {
		t.Error("Expected at least one client resource file")
	}
}

func TestRepoFileContent(t *testing.T) {
	// Test that we can read actual file content
	data, err := Repo.ReadFile("catalogs/all")
	if err != nil {
		t.Fatalf("Failed to read catalogs/all file: %v", err)
	}

	// Should have some content
	if len(data) == 0 {
		t.Error("Expected catalogs/all to have content")
	}
}
