package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mstgnz/goflow/pkg/workflow"
)

func main() {
	// Define command-line flags
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runFile := runCmd.String("file", "", "Path to the workflow file")

	// Parse command-line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Handle commands
	switch os.Args[1] {
	case "run":
		err := runCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
			os.Exit(1)
		}

		if *runFile == "" {
			fmt.Fprintf(os.Stderr, "Error: -file flag is required\n")
			runCmd.Usage()
			os.Exit(1)
		}

		runWorkflow(*runFile)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  goflow run -file <workflow-file>")
}

func runWorkflow(filePath string) {
	// Create a new workflow engine
	engine := workflow.NewEngine()

	// Register default tasks
	engine.RegisterDefaultTasks()

	// Load the workflow
	err := engine.Load(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading workflow: %v\n", err)
		os.Exit(1)
	}

	// Get the workflow name from the file
	workflowName := getWorkflowNameFromFile(filePath)
	if workflowName == "" {
		fmt.Fprintf(os.Stderr, "Error: could not determine workflow name\n")
		os.Exit(1)
	}

	// Run the workflow
	fmt.Printf("Running workflow: %s\n", workflowName)
	state, err := engine.Run(workflowName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running workflow: %v\n", err)
		os.Exit(1)
	}

	// Print the result
	fmt.Printf("Workflow completed with status: %s\n", state.Status)

	// Print the step results
	fmt.Println("Step results:")
	for stepID, result := range state.StepResults {
		fmt.Printf("  %s: %v\n", stepID, result.Success)
	}
}

func getWorkflowNameFromFile(filePath string) string {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	// Parse the JSON
	var workflow struct {
		Name string `json:"name"`
	}

	err = json.Unmarshal(data, &workflow)
	if err != nil {
		return ""
	}

	return workflow.Name
}
