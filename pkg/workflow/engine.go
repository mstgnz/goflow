package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mstgnz/goflow/pkg/models"
	"github.com/mstgnz/goflow/pkg/tasks"
)

// Engine is the core workflow engine
type Engine struct {
	taskRegistry *tasks.Registry
	workflows    map[string]*models.Workflow
	states       map[string]*models.WorkflowState
}

// NewEngine creates a new workflow engine
func NewEngine() *Engine {
	return &Engine{
		taskRegistry: tasks.NewRegistry(),
		workflows:    make(map[string]*models.Workflow),
		states:       make(map[string]*models.WorkflowState),
	}
}

// RegisterTask registers a task with the engine
func (e *Engine) RegisterTask(task tasks.Task) {
	e.taskRegistry.Register(task)
}

// RegisterDefaultTasks registers the default tasks with the engine
func (e *Engine) RegisterDefaultTasks() {
	e.RegisterTask(&tasks.SendEmailTask{})
	e.RegisterTask(&tasks.ProcessPaymentTask{})
	e.RegisterTask(&tasks.PackItemsTask{})
	e.RegisterTask(&tasks.SendShippingNotificationTask{})
	e.RegisterTask(&tasks.ValidateFileTask{})
	e.RegisterTask(&tasks.ProcessFileTask{})
	e.RegisterTask(&tasks.SaveToDatabaseTask{})
}

// Load loads a workflow from a file
func (e *Engine) Load(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow models.Workflow
	if strings.HasSuffix(filePath, ".json") {
		err = json.Unmarshal(data, &workflow)
	} else {
		return fmt.Errorf("unsupported file format: %s", filePath)
	}

	if err != nil {
		return fmt.Errorf("failed to parse workflow file: %w", err)
	}

	if workflow.Name == "" {
		return errors.New("workflow must have a name")
	}

	e.workflows[workflow.Name] = &workflow
	return nil
}

// Run runs a workflow by name
func (e *Engine) Run(workflowName string) (*models.WorkflowState, error) {
	workflow, ok := e.workflows[workflowName]
	if !ok {
		return nil, fmt.Errorf("workflow not found: %s", workflowName)
	}

	// Create a new workflow state
	state := &models.WorkflowState{
		WorkflowName:   workflowName,
		CompletedSteps: []string{},
		StepResults:    make(map[string]models.StepResult),
		StartTime:      time.Now().Unix(),
		Status:         "running",
	}

	// Store the state
	e.states[workflowName] = state

	// Start the workflow execution
	ctx := context.Background()
	err := e.executeWorkflow(ctx, workflow, state)
	if err != nil {
		state.Status = "failed"
		return state, err
	}

	state.Status = "completed"
	state.EndTime = time.Now().Unix()
	return state, nil
}

// executeWorkflow executes a workflow
func (e *Engine) executeWorkflow(ctx context.Context, workflow *models.Workflow, state *models.WorkflowState) error {
	// Find the first step
	if len(workflow.Steps) == 0 {
		return errors.New("workflow has no steps")
	}

	// Start with the first step
	currentStep := workflow.Steps[0]
	state.CurrentStep = currentStep.ID

	// Execute steps until there are no more steps to execute
	for {
		// Execute the current step
		err := e.executeStep(ctx, workflow, currentStep, state)
		if err != nil {
			return fmt.Errorf("failed to execute step %s: %w", currentStep.ID, err)
		}

		// Mark the step as completed
		state.CompletedSteps = append(state.CompletedSteps, currentStep.ID)

		// Find the next step to execute
		nextStepID := e.findNextStep(workflow, currentStep, state)
		if nextStepID == "" {
			// No more steps to execute
			break
		}

		// Find the next step
		var found bool
		for _, step := range workflow.Steps {
			if step.ID == nextStepID {
				currentStep = step
				state.CurrentStep = currentStep.ID
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("step not found: %s", nextStepID)
		}
	}

	return nil
}

// executeStep executes a single step in a workflow
func (e *Engine) executeStep(ctx context.Context, workflow *models.Workflow, step models.Step, state *models.WorkflowState) error {
	// Check if the step has a condition
	if step.Condition != "" {
		// Evaluate the condition
		if !e.evaluateCondition(step.Condition, state) {
			// Condition not met, skip this step
			return nil
		}
	}

	// Get the task
	task, ok := e.taskRegistry.Get(step.Task)
	if !ok {
		return fmt.Errorf("task not found: %s", step.Task)
	}

	// Execute the task
	result, err := task.Execute(ctx, step.Params, state)
	if err != nil {
		// Store the result
		state.StepResults[step.ID] = models.StepResult{
			Success: false,
			Error:   err.Error(),
		}
		return err
	}

	// Store the result
	state.StepResults[step.ID] = models.StepResult{
		Success: true,
		Data:    result,
	}

	return nil
}

// findNextStep finds the next step to execute
func (e *Engine) findNextStep(workflow *models.Workflow, currentStep models.Step, state *models.WorkflowState) string {
	// If there are no next steps, we're done
	if len(currentStep.Next) == 0 {
		return ""
	}

	// Get the result of the current step
	result, ok := state.StepResults[currentStep.ID]
	if !ok || !result.Success {
		// Step failed or no result, don't continue
		return ""
	}

	// Return the first next step
	return currentStep.Next[0]
}

// evaluateCondition evaluates a condition
func (e *Engine) evaluateCondition(condition string, state *models.WorkflowState) bool {
	// Parse the condition (format: "step.field")
	parts := strings.Split(condition, ".")
	if len(parts) != 2 {
		return false
	}

	stepID := parts[0]
	field := parts[1]

	// Get the step result
	result, ok := state.StepResults[stepID]
	if !ok {
		return false
	}

	// Check if the field exists and is true
	value, ok := result.Data[field]
	if !ok {
		return false
	}

	// Convert to bool
	boolValue, ok := value.(bool)
	if !ok {
		return false
	}

	return boolValue
}

// GetState returns the state of a workflow
func (e *Engine) GetState(workflowName string) (*models.WorkflowState, bool) {
	state, ok := e.states[workflowName]
	return state, ok
}
