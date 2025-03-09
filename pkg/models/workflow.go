package models

// Workflow represents a complete workflow definition
type Workflow struct {
	Name  string `json:"name" yaml:"name"`
	Steps []Step `json:"steps" yaml:"steps"`
}

// Step represents a single step in a workflow
type Step struct {
	ID        string            `json:"id" yaml:"id"`
	Task      string            `json:"task" yaml:"task"`
	Next      []string          `json:"next" yaml:"next"`
	Condition string            `json:"condition,omitempty" yaml:"condition,omitempty"`
	Params    map[string]string `json:"params,omitempty" yaml:"params,omitempty"`
}

// WorkflowState represents the current state of a workflow execution
type WorkflowState struct {
	WorkflowName   string                `json:"workflow_name"`
	CurrentStep    string                `json:"current_step"`
	CompletedSteps []string              `json:"completed_steps"`
	StepResults    map[string]StepResult `json:"step_results"`
	StartTime      int64                 `json:"start_time"`
	EndTime        int64                 `json:"end_time,omitempty"`
	Status         string                `json:"status"` // "running", "completed", "failed"
}

// StepResult represents the result of a step execution
type StepResult struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}
