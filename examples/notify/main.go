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

	notifiedWaiters, err := client.Notify("event-id", map[string]string{
		"userId": "testUser",
	})
	if err != nil {
		fmt.Printf("failed to notify event: %v\n", err)
		return
	}
	for _, notifiedResp := range notifiedWaiters {
		fmt.Printf(notifiedResp.Waiter.Url)
	}
}
