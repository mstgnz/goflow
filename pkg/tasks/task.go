package tasks

import (
	"context"

	"github.com/mstgnz/goflow/pkg/models"
)

// Task is the interface that all tasks must implement
type Task interface {
	// Execute runs the task with the given parameters and returns a result
	Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]any, error)
	// Name returns the name of the task
	Name() string
}

// Registry is a registry of all available tasks
type Registry struct {
	tasks map[string]Task
}

// NewRegistry creates a new task registry
func NewRegistry() *Registry {
	return &Registry{
		tasks: make(map[string]Task),
	}
}

// Register registers a task with the registry
func (r *Registry) Register(task Task) {
	r.tasks[task.Name()] = task
}

// Get returns a task by name
func (r *Registry) Get(name string) (Task, bool) {
	task, ok := r.tasks[name]
	return task, ok
}

// List returns a list of all registered task names
func (r *Registry) List() []string {
	var names []string
	for name := range r.tasks {
		names = append(names, name)
	}
	return names
}
