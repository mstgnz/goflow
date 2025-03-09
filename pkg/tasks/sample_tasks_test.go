package tasks

import (
	"context"
	"testing"

	"github.com/mstgnz/goflow/pkg/models"
)

func TestSendEmailTask(t *testing.T) {
	// Create a task
	task := &SendEmailTask{}

	// Verify the task name
	if task.Name() != "send_email" {
		t.Errorf("Expected task name send_email, got %s", task.Name())
	}

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Test with missing template parameter
	_, err := task.Execute(context.Background(), map[string]string{}, state)
	if err == nil {
		t.Error("Expected error for missing template parameter")
	}

	// Test with valid parameters
	params := map[string]string{
		"template": "test_template",
	}
	result, err := task.Execute(context.Background(), params, state)
	if err != nil {
		t.Fatalf("Failed to execute task: %v", err)
	}

	// Verify the result
	if result["sent"] != true {
		t.Errorf("Expected sent true, got %v", result["sent"])
	}

	if result["template"] != "test_template" {
		t.Errorf("Expected template test_template, got %v", result["template"])
	}

	if result["time"] == nil {
		t.Error("Expected time to be set")
	}
}

func TestProcessPaymentTask(t *testing.T) {
	// Create a task
	task := &ProcessPaymentTask{}

	// Verify the task name
	if task.Name() != "process_payment" {
		t.Errorf("Expected task name process_payment, got %s", task.Name())
	}

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Test with missing amount parameter
	_, err := task.Execute(context.Background(), map[string]string{}, state)
	if err == nil {
		t.Error("Expected error for missing amount parameter")
	}

	// Test with valid parameters
	params := map[string]string{
		"amount": "100.00",
	}
	result, err := task.Execute(context.Background(), params, state)
	if err != nil {
		t.Fatalf("Failed to execute task: %v", err)
	}

	// Verify the result
	if result["success"] != true {
		t.Errorf("Expected success true, got %v", result["success"])
	}

	if result["amount"] != "100.00" {
		t.Errorf("Expected amount 100.00, got %v", result["amount"])
	}

	if result["time"] == nil {
		t.Error("Expected time to be set")
	}
}

func TestPackItemsTask(t *testing.T) {
	// Create a task
	task := &PackItemsTask{}

	// Verify the task name
	if task.Name() != "pack_items" {
		t.Errorf("Expected task name pack_items, got %s", task.Name())
	}

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Execute the task
	result, err := task.Execute(context.Background(), map[string]string{}, state)
	if err != nil {
		t.Fatalf("Failed to execute task: %v", err)
	}

	// Verify the result
	if result["packed"] != true {
		t.Errorf("Expected packed true, got %v", result["packed"])
	}

	if result["time"] == nil {
		t.Error("Expected time to be set")
	}
}

func TestSendShippingNotificationTask(t *testing.T) {
	// Create a task
	task := &SendShippingNotificationTask{}

	// Verify the task name
	if task.Name() != "send_shipping_notification" {
		t.Errorf("Expected task name send_shipping_notification, got %s", task.Name())
	}

	// Create a workflow state
	state := &models.WorkflowState{
		WorkflowName:   "test_workflow",
		CurrentStep:    "step1",
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		Status:         "running",
	}

	// Execute the task
	result, err := task.Execute(context.Background(), map[string]string{}, state)
	if err != nil {
		t.Fatalf("Failed to execute task: %v", err)
	}

	// Verify the result
	if result["sent"] != true {
		t.Errorf("Expected sent true, got %v", result["sent"])
	}

	if result["time"] == nil {
		t.Error("Expected time to be set")
	}
}
