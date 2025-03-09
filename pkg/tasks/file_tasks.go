package tasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mstgnz/goflow/pkg/models"
)

// ValidateFileTask validates a file
type ValidateFileTask struct{}

func (t *ValidateFileTask) Name() string {
	return "validate_file"
}

func (t *ValidateFileTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]interface{}, error) {
	filePath, ok := params["file_path"]
	if !ok {
		return nil, errors.New("file_path parameter is required")
	}

	// In a real implementation, this would validate the file
	fmt.Printf("Validating file: %s\n", filePath)

	// Simulate some work
	time.Sleep(1 * time.Second)

	// Simulate validation (in a real implementation, this could fail)
	valid := true

	return map[string]interface{}{
		"valid":     valid,
		"file_path": filePath,
		"time":      time.Now().Format(time.RFC3339),
	}, nil
}

// ProcessFileTask processes a file
type ProcessFileTask struct{}

func (t *ProcessFileTask) Name() string {
	return "process_file"
}

func (t *ProcessFileTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]interface{}, error) {
	filePath, ok := params["file_path"]
	if !ok {
		return nil, errors.New("file_path parameter is required")
	}

	// In a real implementation, this would process the file
	fmt.Printf("Processing file: %s\n", filePath)

	// Simulate some work
	time.Sleep(2 * time.Second)

	// Simulate processing (in a real implementation, this could fail)
	processed := true

	return map[string]interface{}{
		"processed": processed,
		"file_path": filePath,
		"records":   100, // Simulated number of records processed
		"time":      time.Now().Format(time.RFC3339),
	}, nil
}

// SaveToDatabaseTask saves data to a database
type SaveToDatabaseTask struct{}

func (t *SaveToDatabaseTask) Name() string {
	return "save_to_database"
}

func (t *SaveToDatabaseTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]interface{}, error) {
	// In a real implementation, this would save data to a database
	fmt.Println("Saving data to database")

	// Get the number of records from the previous step
	var records int
	if result, ok := state.StepResults["process"]; ok {
		if recordsVal, ok := result.Data["records"]; ok {
			if recordsInt, ok := recordsVal.(int); ok {
				records = recordsInt
			}
		}
	}

	// Simulate some work
	time.Sleep(1500 * time.Millisecond)

	return map[string]interface{}{
		"saved":   true,
		"records": records,
		"time":    time.Now().Format(time.RFC3339),
	}, nil
}
