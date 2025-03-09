package tasks

import (
	"context"
	"testing"

	"github.com/mstgnz/goflow/pkg/models"
)

// MockTask is a mock task for testing
type MockTask struct {
	name     string
	executed bool
	params   map[string]string
}

func (t *MockTask) Name() string {
	return t.name
}

func (t *MockTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]any, error) {
	t.executed = true
	t.params = params
	return map[string]any{
		"success": true,
		"mock":    "data",
	}, nil
}

func TestTaskRegistry(t *testing.T) {
	// Create a new registry
	registry := NewRegistry()
	if registry == nil {
		t.Fatal("Failed to create registry")
	}

	// Register a task
	task1 := &MockTask{name: "task1"}
	registry.Register(task1)

	// Register another task
	task2 := &MockTask{name: "task2"}
	registry.Register(task2)

	// Get a task
	retrievedTask, ok := registry.Get("task1")
	if !ok {
		t.Fatal("Failed to get task1")
	}

	// Verify the task
	if retrievedTask.Name() != "task1" {
		t.Errorf("Expected task name task1, got %s", retrievedTask.Name())
	}

	// Get a non-existent task
	_, ok = registry.Get("non-existent")
	if ok {
		t.Error("Expected non-existent task to not be found")
	}

	// List tasks
	tasks := registry.List()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Check if the tasks are in the list
	found1 := false
	found2 := false
	for _, name := range tasks {
		if name == "task1" {
			found1 = true
		}
		if name == "task2" {
			found2 = true
		}
	}

	if !found1 {
		t.Error("Expected task1 to be in the list")
	}

	if !found2 {
		t.Error("Expected task2 to be in the list")
	}
}

func TestTaskExecution(t *testing.T) {
	// Create a mock task
	task := &MockTask{name: "test_task"}

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Execute the task
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	result, err := task.Execute(context.Background(), params, state)
	if err != nil {
		t.Fatalf("Failed to execute task: %v", err)
	}

	// Verify the task was executed
	if !task.executed {
		t.Error("Expected task to be executed")
	}

	// Verify the parameters were passed correctly
	if task.params["param1"] != "value1" {
		t.Errorf("Expected param1 value1, got %s", task.params["param1"])
	}

	if task.params["param2"] != "value2" {
		t.Errorf("Expected param2 value2, got %s", task.params["param2"])
	}

	// Verify the result
	if result["success"] != true {
		t.Errorf("Expected success true, got %v", result["success"])
	}

	if result["mock"] != "data" {
		t.Errorf("Expected mock data, got %v", result["mock"])
	}
}
