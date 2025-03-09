package workflow

import (
	"os"
	"testing"
)

func TestIntegrationOrderProcess(t *testing.T) {
	// Create a new engine
	engine := NewEngine()

	// Register the default tasks
	engine.RegisterDefaultTasks()

	// Create a temporary workflow file
	workflowJSON := `{
		"name": "order_process",
		"steps": [
			{
				"id": "payment",
				"task": "process_payment",
				"next": ["prepare_order"],
				"params": {
					"amount": "100.00"
				}
			},
			{
				"id": "prepare_order",
				"task": "pack_items",
				"next": ["ship_order"],
				"condition": "payment.success"
			},
			{
				"id": "ship_order",
				"task": "send_shipping_notification",
				"next": ["thank_you"]
			},
			{
				"id": "thank_you",
				"task": "send_email",
				"params": {
					"template": "thank_you"
				}
			}
		]
	}`

	tmpfile, err := os.CreateTemp("", "order-process-*.json")
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

	// Run the workflow
	state, err := engine.Run("order_process")
	if err != nil {
		t.Fatalf("Failed to run workflow: %v", err)
	}

	// Verify the workflow completed successfully
	if state.Status != "completed" {
		t.Errorf("Expected status completed, got %s", state.Status)
	}

	// Verify all steps were completed
	expectedSteps := []string{"payment", "prepare_order", "ship_order", "thank_you"}
	if len(state.CompletedSteps) != len(expectedSteps) {
		t.Errorf("Expected %d completed steps, got %d", len(expectedSteps), len(state.CompletedSteps))
	}

	for i, step := range expectedSteps {
		if i < len(state.CompletedSteps) && state.CompletedSteps[i] != step {
			t.Errorf("Expected step %s at position %d, got %s", step, i, state.CompletedSteps[i])
		}
	}

	// Verify all steps have results
	for _, step := range expectedSteps {
		result, ok := state.StepResults[step]
		if !ok {
			t.Errorf("Expected result for step %s, but not found", step)
			continue
		}

		if !result.Success {
			t.Errorf("Expected step %s to succeed, but it failed: %s", step, result.Error)
		}
	}
}

func TestIntegrationFileProcessing(t *testing.T) {
	// Create a new engine
	engine := NewEngine()

	// Register the default tasks
	engine.RegisterDefaultTasks()

	// Create a temporary workflow file
	workflowJSON := `{
		"name": "file_processing",
		"steps": [
			{
				"id": "validate",
				"task": "validate_file",
				"next": ["process"],
				"params": {
					"file_path": "/path/to/data.csv"
				}
			},
			{
				"id": "process",
				"task": "process_file",
				"next": ["save"],
				"condition": "validate.valid",
				"params": {
					"file_path": "/path/to/data.csv"
				}
			},
			{
				"id": "save",
				"task": "save_to_database",
				"next": ["notify"],
				"condition": "process.processed"
			},
			{
				"id": "notify",
				"task": "send_email",
				"params": {
					"template": "processing_complete"
				}
			}
		]
	}`

	tmpfile, err := os.CreateTemp("", "file-processing-*.json")
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

	// Run the workflow
	state, err := engine.Run("file_processing")
	if err != nil {
		t.Fatalf("Failed to run workflow: %v", err)
	}

	// Verify the workflow completed successfully
	if state.Status != "completed" {
		t.Errorf("Expected status completed, got %s", state.Status)
	}

	// Verify all steps were completed
	expectedSteps := []string{"validate", "process", "save", "notify"}
	if len(state.CompletedSteps) != len(expectedSteps) {
		t.Errorf("Expected %d completed steps, got %d", len(expectedSteps), len(state.CompletedSteps))
	}

	for i, step := range expectedSteps {
		if i < len(state.CompletedSteps) && state.CompletedSteps[i] != step {
			t.Errorf("Expected step %s at position %d, got %s", step, i, state.CompletedSteps[i])
		}
	}

	// Verify all steps have results
	for _, step := range expectedSteps {
		result, ok := state.StepResults[step]
		if !ok {
			t.Errorf("Expected result for step %s, but not found", step)
			continue
		}

		if !result.Success {
			t.Errorf("Expected step %s to succeed, but it failed: %s", step, result.Error)
		}
	}

	// Verify the records were processed and saved
	processResult, ok := state.StepResults["process"]
	if !ok {
		t.Fatal("Expected result for process step, but not found")
	}

	records, ok := processResult.Data["records"]
	if !ok {
		t.Fatal("Expected records in process result, but not found")
	}

	// Check that records is a number (could be int or float64)
	switch records.(type) {
	case int, float64:
		// Valid type
	default:
		t.Errorf("Expected records to be a number, got %T", records)
	}

	saveResult, ok := state.StepResults["save"]
	if !ok {
		t.Fatal("Expected result for save step, but not found")
	}

	savedRecords, ok := saveResult.Data["records"]
	if !ok {
		t.Fatal("Expected records in save result, but not found")
	}

	// Check that savedRecords is a number (could be int or float64)
	switch savedRecords.(type) {
	case int, float64:
		// Valid type
	default:
		t.Errorf("Expected saved records to be a number, got %T", savedRecords)
	}
}
