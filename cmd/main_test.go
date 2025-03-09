package main

import (
	"os"
	"testing"
)

func TestGetWorkflowNameFromFile(t *testing.T) {
	// Create a temporary workflow file
	workflowJSON := `{
		"name": "test_workflow",
		"steps": []
	}`

	tmpfile, err := os.CreateTemp("", "workflow-*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(workflowJSON)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Get the workflow name
	name := getWorkflowNameFromFile(tmpfile.Name())
	if name != "test_workflow" {
		t.Errorf("Expected workflow name test_workflow, got %s", name)
	}

	// Test with a non-existent file
	name = getWorkflowNameFromFile("non-existent-file.json")
	if name != "" {
		t.Errorf("Expected empty workflow name for non-existent file, got %s", name)
	}

	// Test with an invalid JSON file
	invalidJSON := `{
		"name": "test_workflow",
		"steps": [
	}`

	tmpfile2, err := os.CreateTemp("", "invalid-*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile2.Name())

	if _, err := tmpfile2.Write([]byte(invalidJSON)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile2.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	name = getWorkflowNameFromFile(tmpfile2.Name())
	if name != "" {
		t.Errorf("Expected empty workflow name for invalid JSON, got %s", name)
	}
}
