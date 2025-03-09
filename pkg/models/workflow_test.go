package models

import (
	"encoding/json"
	"testing"
)

func TestWorkflowSerialization(t *testing.T) {
	// Create a sample workflow
	workflow := Workflow{
		Name: "test_workflow",
		Steps: []Step{
			{
				ID:   "step1",
				Task: "task1",
				Next: []string{"step2"},
				Params: map[string]string{
					"param1": "value1",
				},
			},
			{
				ID:        "step2",
				Task:      "task2",
				Next:      []string{"step3"},
				Condition: "step1.success",
			},
			{
				ID:   "step3",
				Task: "task3",
			},
		},
	}

	// Serialize to JSON
	data, err := json.Marshal(workflow)
	if err != nil {
		t.Fatalf("Failed to marshal workflow: %v", err)
	}

	// Deserialize from JSON
	var deserializedWorkflow Workflow
	err = json.Unmarshal(data, &deserializedWorkflow)
	if err != nil {
		t.Fatalf("Failed to unmarshal workflow: %v", err)
	}

	// Verify the deserialized workflow
	if deserializedWorkflow.Name != workflow.Name {
		t.Errorf("Expected workflow name %s, got %s", workflow.Name, deserializedWorkflow.Name)
	}

	if len(deserializedWorkflow.Steps) != len(workflow.Steps) {
		t.Errorf("Expected %d steps, got %d", len(workflow.Steps), len(deserializedWorkflow.Steps))
	}

	// Check the first step
	if deserializedWorkflow.Steps[0].ID != "step1" {
		t.Errorf("Expected step ID step1, got %s", deserializedWorkflow.Steps[0].ID)
	}

	if deserializedWorkflow.Steps[0].Task != "task1" {
		t.Errorf("Expected task task1, got %s", deserializedWorkflow.Steps[0].Task)
	}

	if len(deserializedWorkflow.Steps[0].Next) != 1 || deserializedWorkflow.Steps[0].Next[0] != "step2" {
		t.Errorf("Expected next step step2, got %v", deserializedWorkflow.Steps[0].Next)
	}

	if deserializedWorkflow.Steps[0].Params["param1"] != "value1" {
		t.Errorf("Expected param1 value1, got %s", deserializedWorkflow.Steps[0].Params["param1"])
	}

	// Check the second step
	if deserializedWorkflow.Steps[1].Condition != "step1.success" {
		t.Errorf("Expected condition step1.success, got %s", deserializedWorkflow.Steps[1].Condition)
	}
}

func TestWorkflowState(t *testing.T) {
	// Create a sample workflow state
	state := WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]StepResult),
		StartTime:      1234567890,
		Status:         "running",
	}

	// Add a step result
	state.StepResults["step1"] = StepResult{
		Success: true,
		Data: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}

	// Serialize to JSON
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal workflow state: %v", err)
	}

	// Deserialize from JSON
	var deserializedState WorkflowState
	err = json.Unmarshal(data, &deserializedState)
	if err != nil {
		t.Fatalf("Failed to unmarshal workflow state: %v", err)
	}

	// Verify the deserialized state
	if deserializedState.WorkflowName != state.WorkflowName {
		t.Errorf("Expected workflow name %s, got %s", state.WorkflowName, deserializedState.WorkflowName)
	}

	if deserializedState.CurrentStep != state.CurrentStep {
		t.Errorf("Expected current step %s, got %s", state.CurrentStep, deserializedState.CurrentStep)
	}

	if deserializedState.Status != state.Status {
		t.Errorf("Expected status %s, got %s", state.Status, deserializedState.Status)
	}

	// Check the step result
	result, ok := deserializedState.StepResults["step1"]
	if !ok {
		t.Fatalf("Expected step result for step1, but not found")
	}

	if !result.Success {
		t.Errorf("Expected success true, got false")
	}

	if result.Data["key1"] != "value1" {
		t.Errorf("Expected key1 value1, got %v", result.Data["key1"])
	}

	if result.Data["key2"] != float64(123) {
		t.Errorf("Expected key2 123, got %v", result.Data["key2"])
	}
}
