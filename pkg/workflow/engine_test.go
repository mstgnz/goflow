package workflow

import (
	"context"
	"os"
	"testing"

	"github.com/mstgnz/goflow/pkg/models"
)

// MockTask is a mock task for testing
type MockTask struct {
	name     string
	executed bool
	params   map[string]string
	result   map[string]interface{}
	err      error
}

func (t *MockTask) Name() string {
	return t.name
}

func (t *MockTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]interface{}, error) {
	t.executed = true
	t.params = params
	return t.result, t.err
}

func TestEngineCreation(t *testing.T) {
	// Create a new engine
	engine := NewEngine()
	if engine == nil {
		t.Fatal("Failed to create engine")
	}

	// Verify the engine has an empty task registry
	if engine.taskRegistry == nil {
		t.Error("Expected task registry to be initialized")
	}

	// Verify the engine has an empty workflow map
	if engine.workflows == nil {
		t.Error("Expected workflows map to be initialized")
	}

	// Verify the engine has an empty state map
	if engine.states == nil {
		t.Error("Expected states map to be initialized")
	}
}

func TestTaskRegistration(t *testing.T) {
	// Create a new engine
	engine := NewEngine()

	// Register a task
	task := &MockTask{name: "mock_task", result: map[string]interface{}{"success": true}}
	engine.RegisterTask(task)

	// Verify the task is registered
	registeredTask, ok := engine.taskRegistry.Get("mock_task")
	if !ok {
		t.Fatal("Failed to get registered task")
	}

	// Verify the task name
	if registeredTask.Name() != "mock_task" {
		t.Errorf("Expected task name mock_task, got %s", registeredTask.Name())
	}
}

func TestWorkflowLoading(t *testing.T) {
	// Create a new engine
	engine := NewEngine()

	// Create a temporary workflow file
	workflowJSON := `{
		"name": "test_workflow",
		"steps": [
			{
				"id": "step1",
				"task": "task1",
				"next": ["step2"]
			},
			{
				"id": "step2",
				"task": "task2"
			}
		]
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

	// Load the workflow
	err = engine.Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load workflow: %v", err)
	}

	// Verify the workflow is loaded
	workflow, ok := engine.workflows["test_workflow"]
	if !ok {
		t.Fatal("Failed to get loaded workflow")
	}

	// Verify the workflow name
	if workflow.Name != "test_workflow" {
		t.Errorf("Expected workflow name test_workflow, got %s", workflow.Name)
	}

	// Verify the workflow steps
	if len(workflow.Steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(workflow.Steps))
	}

	// Verify the first step
	if workflow.Steps[0].ID != "step1" {
		t.Errorf("Expected step ID step1, got %s", workflow.Steps[0].ID)
	}

	if workflow.Steps[0].Task != "task1" {
		t.Errorf("Expected task task1, got %s", workflow.Steps[0].Task)
	}

	if len(workflow.Steps[0].Next) != 1 || workflow.Steps[0].Next[0] != "step2" {
		t.Errorf("Expected next step step2, got %v", workflow.Steps[0].Next)
	}

	// Verify the second step
	if workflow.Steps[1].ID != "step2" {
		t.Errorf("Expected step ID step2, got %s", workflow.Steps[1].ID)
	}

	if workflow.Steps[1].Task != "task2" {
		t.Errorf("Expected task task2, got %s", workflow.Steps[1].Task)
	}
}

func TestWorkflowExecution(t *testing.T) {
	// Create a new engine
	engine := NewEngine()

	// Register mock tasks
	task1 := &MockTask{name: "task1", result: map[string]interface{}{"success": true}}
	task2 := &MockTask{name: "task2", result: map[string]interface{}{"success": true}}
	engine.RegisterTask(task1)
	engine.RegisterTask(task2)

	// Create a workflow
	workflow := &models.Workflow{
		Name: "test_workflow",
		Steps: []models.Step{
			{
				ID:   "step1",
				Task: "task1",
				Next: []string{"step2"},
			},
			{
				ID:        "step2",
				Task:      "task2",
				Condition: "step1.success",
			},
		},
	}

	// Add the workflow to the engine
	engine.workflows["test_workflow"] = workflow

	// Run the workflow
	state, err := engine.Run("test_workflow")
	if err != nil {
		t.Fatalf("Failed to run workflow: %v", err)
	}

	// Verify the workflow state
	if state.WorkflowName != "test_workflow" {
		t.Errorf("Expected workflow name test_workflow, got %s", state.WorkflowName)
	}

	if state.Status != "completed" {
		t.Errorf("Expected status completed, got %s", state.Status)
	}

	// Verify the tasks were executed
	if !task1.executed {
		t.Error("Expected task1 to be executed")
	}

	if !task2.executed {
		t.Error("Expected task2 to be executed")
	}

	// Verify the step results
	if len(state.StepResults) != 2 {
		t.Errorf("Expected 2 step results, got %d", len(state.StepResults))
	}

	// Verify the completed steps
	if len(state.CompletedSteps) != 2 {
		t.Errorf("Expected 2 completed steps, got %d", len(state.CompletedSteps))
	}

	if state.CompletedSteps[0] != "step1" {
		t.Errorf("Expected first completed step to be step1, got %s", state.CompletedSteps[0])
	}

	if state.CompletedSteps[1] != "step2" {
		t.Errorf("Expected second completed step to be step2, got %s", state.CompletedSteps[1])
	}
}

func TestConditionEvaluation(t *testing.T) {
	// Create a new engine
	engine := NewEngine()

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Add a step result with success=true
	state.StepResults["step1"] = models.StepResult{
		Success: true,
		Data: map[string]interface{}{
			"success": true,
			"value":   123,
		},
	}

	// Test a condition that should evaluate to true
	condition := "step1.success"
	result := engine.evaluateCondition(condition, state)
	if !result {
		t.Errorf("Expected condition %s to evaluate to true", condition)
	}

	// Test a condition that should evaluate to false
	condition = "step1.failure"
	result = engine.evaluateCondition(condition, state)
	if result {
		t.Errorf("Expected condition %s to evaluate to false", condition)
	}

	// Test a condition with a non-existent step
	condition = "non_existent.success"
	result = engine.evaluateCondition(condition, state)
	if result {
		t.Errorf("Expected condition %s to evaluate to false", condition)
	}

	// Test a condition with an invalid format
	condition = "invalid_condition"
	result = engine.evaluateCondition(condition, state)
	if result {
		t.Errorf("Expected condition %s to evaluate to false", condition)
	}
}

func TestGetState(t *testing.T) {
	// Create a new engine
	engine := NewEngine()

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Add the state to the engine
	engine.states["test_workflow"] = state

	// Get the state
	retrievedState, ok := engine.GetState("test_workflow")
	if !ok {
		t.Fatal("Failed to get workflow state")
	}

	// Verify the state
	if retrievedState.WorkflowName != "test_workflow" {
		t.Errorf("Expected workflow name test_workflow, got %s", retrievedState.WorkflowName)
	}

	if retrievedState.CurrentStep != "step1" {
		t.Errorf("Expected current step step1, got %s", retrievedState.CurrentStep)
	}

	if retrievedState.Status != "running" {
		t.Errorf("Expected status running, got %s", retrievedState.Status)
	}

	// Get a non-existent state
	_, ok = engine.GetState("non_existent")
	if ok {
		t.Error("Expected non-existent state to not be found")
	}
}
