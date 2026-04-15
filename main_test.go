package main

import (
	"os"
	"testing"
)

func TestSaveAndLoadManifest(t *testing.T) {
	testFile := "test-docker-deps.json"
	defer func() {
		_ = os.Remove(testFile)
	}()

	m := &Manifest{
		Name:    "test-app",
		Version: "1.0.0",
		Dependencies: map[string]string{
			"nginx": "1.25.1",
		},
	}

	if err := saveManifest(m, testFile); err != nil {
		t.Fatalf("saveManifest failed: %v", err)
	}

	loaded, err := loadManifest(testFile)
	if err != nil {
		t.Fatalf("loadManifest failed: %v", err)
	}

	if loaded.Name != m.Name {
		t.Errorf("Expected name %s, got %s", m.Name, loaded.Name)
	}
	if loaded.Dependencies["nginx"] != "1.25.1" {
		t.Errorf("Expected version 1.25.1, got %s", loaded.Dependencies["nginx"])
	}
}

func TestLoadNonExistentManifest(t *testing.T) {
	m, err := loadManifest("non-existent.json")
	if err != nil {
		t.Fatalf("loadManifest should not fail for non-existent file, got %v", err)
	}
	if m.Dependencies == nil {
		t.Fatal("Expected non-nil Dependencies map")
	}
}

func TestLoadInvalidManifest(t *testing.T) {
	testFile := "invalid.json"
	_ = os.WriteFile(testFile, []byte("{invalid json}"), 0644)
	defer func() {
		_ = os.Remove(testFile)
	}()

	_, err := loadManifest(testFile)
	if err == nil {
		t.Fatal("Expected error loading invalid JSON")
	}
}

func TestManifestBusinessLogic(t *testing.T) {
	m := &Manifest{Dependencies: make(map[string]string)}

	// Test SetProjectInfo
	m.SetProjectInfo("my-app", "2.0.0")
	if m.Name != "my-app" || m.Version != "2.0.0" {
		t.Errorf("SetProjectInfo failed")
	}

	// Test AddDependency
	m.AddDependency("redis", "7.0")
	if m.Dependencies["redis"] != "7.0" {
		t.Errorf("AddDependency failed")
	}

	// Test GetDependency
	ver, ok := m.GetDependency("redis")
	if !ok || ver != "7.0" {
		t.Errorf("GetDependency failed")
	}

	// Test UpdateDependency
	err := m.UpdateDependency("redis", "7.2")
	if err != nil {
		t.Errorf("UpdateDependency failed: %v", err)
	}
	if m.Dependencies["redis"] != "7.2" {
		t.Errorf("UpdateDependency version mismatch")
	}

	// Test UpdateDependency non-existent
	err = m.UpdateDependency("mysql", "8.0")
	if err == nil {
		t.Errorf("UpdateDependency should fail for non-existent dep")
	}

	// Test RemoveDependency
	err = m.RemoveDependency("redis")
	if err != nil {
		t.Errorf("RemoveDependency failed: %v", err)
	}
	if _, ok := m.Dependencies["redis"]; ok {
		t.Errorf("RemoveDependency did not remove the key")
	}

	// Test RemoveDependency non-existent
	err = m.RemoveDependency("redis")
	if err == nil {
		t.Errorf("RemoveDependency should fail for non-existent dep")
	}
}
