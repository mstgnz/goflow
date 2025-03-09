package tasks

import (
	"context"
	"testing"

	"github.com/mstgnz/goflow/pkg/models"
)

func TestValidateFileTask(t *testing.T) {
	// Create a task
	task := &ValidateFileTask{}

	// Verify the task name
	if task.Name() != "validate_file" {
		t.Errorf("Expected task name validate_file, got %s", task.Name())
	}

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Test with missing file_path parameter
	_, err := task.Execute(context.Background(), map[string]string{}, state)
	if err == nil {
		t.Error("Expected error for missing file_path parameter")
	}

	// Test with valid parameters
	params := map[string]string{
		"file_path": "/path/to/file.csv",
	}
	result, err := task.Execute(context.Background(), params, state)
	if err != nil {
		t.Fatalf("Failed to execute task: %v", err)
	}

	// Verify the result
	if result["valid"] != true {
		t.Errorf("Expected valid true, got %v", result["valid"])
	}

	if result["file_path"] != "/path/to/file.csv" {
		t.Errorf("Expected file_path /path/to/file.csv, got %v", result["file_path"])
	}

	if result["time"] == nil {
		t.Error("Expected time to be set")
	}
}

func TestProcessFileTask(t *testing.T) {
	// Create a task
	task := &ProcessFileTask{}

	// Verify the task name
	if task.Name() != "process_file" {
		t.Errorf("Expected task name process_file, got %s", task.Name())
	}

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Test with missing file_path parameter
	_, err := task.Execute(context.Background(), map[string]string{}, state)
	if err == nil {
		t.Error("Expected error for missing file_path parameter")
	}

	// Test with valid parameters
	params := map[string]string{
		"file_path": "/path/to/file.csv",
	}
	result, err := task.Execute(context.Background(), params, state)
	if err != nil {
		t.Fatalf("Failed to execute task: %v", err)
	}

	// Verify the result
	if result["processed"] != true {
		t.Errorf("Expected processed true, got %v", result["processed"])
	}

	if result["file_path"] != "/path/to/file.csv" {
		t.Errorf("Expected file_path /path/to/file.csv, got %v", result["file_path"])
	}

	if result["records"] != 100 {
		t.Errorf("Expected records 100, got %v", result["records"])
	}

	if result["time"] == nil {
		t.Error("Expected time to be set")
	}
}

func TestSaveToDatabaseTask(t *testing.T) {
	// Create a task
	task := &SaveToDatabaseTask{}

	// Verify the task name
	if task.Name() != "save_to_database" {
		t.Errorf("Expected task name save_to_database, got %s", task.Name())
	}

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Add a previous step result with records
	state.StepResults["process"] = models.StepResult{
		Success: true,
		Data: map[string]interface{}{
			"records": 100,
		},
	}

	// Execute the task
	result, err := task.Execute(context.Background(), map[string]string{}, state)
	if err != nil {
		t.Fatalf("Failed to execute task: %v", err)
	}

	// Verify the result
	if result["saved"] != true {
		t.Errorf("Expected saved true, got %v", result["saved"])
	}

	if result["records"] != 100 {
		t.Errorf("Expected records 100, got %v", result["records"])
	}

	if result["time"] == nil {
		t.Error("Expected time to be set")
	}
}
