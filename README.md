# goflow

Light Workflow Engine

The idea of a **Light Workflow Engine with Go** aims to provide a simpler and lighter alternative compared to large-scale solutions like **n8n** or **Temporal**. Such a system ensures that specific tasks are executed in a specific order, based on specific conditions. It can be used **especially for managing event-based or schedule-based workflows**.

---

## **Example Scenarios:**

1. **Email Verification Flow:**

   - User registers.
   - Email verification code is sent.
   - If the user confirms, the account is activated, otherwise it times out.

2. **File Processing Flow:**

   - User uploads a CSV file.
   - System processes the file, checks for errors.
   - If the file is valid, it is saved to the database.
   - User is notified about the result.

3. **Notification Chain:**
   - User makes a payment.
   - If the payment is successful, the order is prepared.
   - When the order is ready, the shipping process begins.
   - After delivery, a thank you email is sent to the user.

---

## **How It Works?**

A light **workflow engine** is a system that **runs specific tasks in a specified order**. A workflow consists of the following components:

- **Tasks:** Steps in the workflow. For example: `send_email`, `validate_data`, `upload_to_s3`.
- **Conditions:** Rules that determine when a step will run. For example: `continue if the previous task was successful`.
- **Scheduler:** Ensures tasks run at a specific time or when an event occurs.
- **State Management:** Tracks which tasks have been completed and at what stage they are.

---

## **Example Usage:**

### **1. Workflow Definition**

A **JSON/YAML based workflow** definition could be:

```json
{
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
}
```

### **2. Running the Workflow Engine**

It can be run with a terminal command:

```bash
goflow run -file order_process.json
```

Or it can be triggered programmatically using the Go API:

```go
engine := workflow.NewEngine()
engine.RegisterDefaultTasks()
engine.Load("order_process.json")
engine.Run("order_process")
```

---

## **Features:**

✔ **Lightweight:** Can manage JSON/YAML based workflows without needing large and complex systems.  
✔ **Modular:** New tasks can be easily added.  
✔ **Scheduled Operations:** Tasks can be scheduled to run at specific times.  
✔ **State Management:** Can track which tasks are at which stage.

---

## **Installation and Usage**

### **Installation**

```bash
# Clone the project
git clone https://github.com/mstgnz/goflow.git
cd goflow

# Build the project
make build

# Or run directly
make run
```

### **Adding Your Own Tasks**

To add your own tasks, create a new file in the `pkg/tasks` directory and define a structure that implements the `Task` interface:

```go
package tasks

import (
	"context"
	"github.com/mstgnz/goflow/pkg/models"
)

// MyCustomTask is a custom task
type MyCustomTask struct{}

func (t *MyCustomTask) Name() string {
	return "my_custom_task"
}

func (t *MyCustomTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]interface{}, error) {
	// Implement the functionality of the task here
	return map[string]interface{}{
		"success": true,
	}, nil
}
```

Then, register this task with the workflow engine:

```go
engine := workflow.NewEngine()
engine.RegisterTask(&tasks.MyCustomTask{})
```

### **Project Structure**

```
goflow/
├── cmd/
│   └── main.go           # Main application entry point
├── pkg/
│   ├── models/           # Data models
│   │   └── workflow.go   # Workflow and step models
│   ├── tasks/            # Task definitions
│   │   ├── task.go       # Task interface
│   │   └── sample_tasks.go # Example tasks
│   └── workflow/         # Workflow engine
│       └── engine.go     # Main workflow engine
├── examples/             # Example workflow definitions
│   └── order_process.json # Example order process
├── Makefile              # Build and run commands
└── README.md             # Project documentation
```

---
