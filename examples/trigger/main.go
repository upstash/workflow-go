package main

import (
	"fmt"
	workflow "workflow-go"
)

func main() {
	client, err := workflow.NewClient("<QSTASH_TOKEN>")
	if err != nil {
		fmt.Printf("failed to create client: %v\n", err)
		return
	}

	runID, err := client.Trigger(workflow.TriggerOptions{
		Url: "https://your-workflow-endpoint.com/api/process",
		Body: map[string]interface{}{
			"user_id":    12345,
			"action":     "purchase",
			"items":      []string{"item1", "item2", "item3"},
			"total_cost": 99.99,
			"timestamp":  "2025-04-16T18:05:24+03:00",
		},
		Retries:        workflow.Retry(2),
		FlowControlKey: "my-flow-control-key",
		Rate:           workflow.Rate(100),
		Parallelism:    workflow.Parallelism(100),
	})
	if err != nil {
		fmt.Printf("failed to trigger workflow run: %v\n", err)
		return
	}
	fmt.Println(runID)
}
