package main

import (
	"fmt"
	"github.com/upstash/workflow-go"
)

func main() {
	client, err := workflow.NewClient("<QSTASH_TOKEN>")
	if err != nil {
		fmt.Printf("failed to create client: %v\n", err)
		return
	}

	// Cancel a single workflow run
	err = client.Cancel("workflow-run-id")
	if err != nil {
		fmt.Printf("failed to cancel workflow run: %v\n", err)
		return
	}

	// Cancel a set of workflow runs
	canceled, err := client.CancelMany([]string{"run-a", "run-b", "run-c"})
	if err != nil {
		fmt.Printf("failed to cancel workflow runs: %v\n", err)
		return
	}
	fmt.Printf("canceled %d workflow runs\n", canceled)

	// Cancel all workflow runs
	canceled, err = client.CancelAll()
	if err != nil {
		fmt.Printf("failed to cancel all workflow runs: %v\n", err)
		return
	}
	fmt.Printf("canceled %d workflow runs", canceled)

	// Cancel all workflow runs that start with a specific URL
	canceled, err = client.CancelWithFilter(workflow.CancelFilter{UrlStartingWith: "http://app-domian.com"})
	if err != nil {
		fmt.Printf("failed to cancel workflow runs: %v\n", err)
		return
	}
	fmt.Printf("canceled %d workflow runs", canceled)
}
