package tasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mstgnz/goflow/pkg/models"
)

// SendEmailTask sends an email
type SendEmailTask struct{}

func (t *SendEmailTask) Name() string {
	return "send_email"
}

func (t *SendEmailTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]interface{}, error) {
	template, ok := params["template"]
	if !ok {
		return nil, errors.New("template parameter is required")
	}

	// In a real implementation, this would send an actual email
	fmt.Printf("Sending email with template: %s\n", template)

	// Simulate some work
	time.Sleep(500 * time.Millisecond)

	return map[string]interface{}{
		"sent":     true,
		"template": template,
		"time":     time.Now().Format(time.RFC3339),
	}, nil
}

// ProcessPaymentTask processes a payment
type ProcessPaymentTask struct{}

func (t *ProcessPaymentTask) Name() string {
	return "process_payment"
}

func (t *ProcessPaymentTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]interface{}, error) {
	amount, ok := params["amount"]
	if !ok {
		return nil, errors.New("amount parameter is required")
	}

	// In a real implementation, this would process an actual payment
	fmt.Printf("Processing payment of amount: %s\n", amount)

	// Simulate some work
	time.Sleep(1 * time.Second)

	// Simulate success (in a real implementation, this could fail)
	success := true

	return map[string]interface{}{
		"success": success,
		"amount":  amount,
		"time":    time.Now().Format(time.RFC3339),
	}, nil
}

// PackItemsTask packs items for an order
type PackItemsTask struct{}

func (t *PackItemsTask) Name() string {
	return "pack_items"
}

func (t *PackItemsTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]interface{}, error) {
	// In a real implementation, this would interact with an inventory system
	fmt.Println("Packing items for order")

	// Simulate some work
	time.Sleep(1500 * time.Millisecond)

	return map[string]interface{}{
		"packed": true,
		"time":   time.Now().Format(time.RFC3339),
	}, nil
}

// SendShippingNotificationTask sends a shipping notification
type SendShippingNotificationTask struct{}

func (t *SendShippingNotificationTask) Name() string {
	return "send_shipping_notification"
}

func (t *SendShippingNotificationTask) Execute(ctx context.Context, params map[string]string, state *models.WorkflowState) (map[string]interface{}, error) {
	// In a real implementation, this would send an actual notification
	fmt.Println("Sending shipping notification")

	// Simulate some work
	time.Sleep(500 * time.Millisecond)

	return map[string]interface{}{
		"sent": true,
		"time": time.Now().Format(time.RFC3339),
	}, nil
}
