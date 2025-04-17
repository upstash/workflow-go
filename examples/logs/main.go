package main

import (
	"fmt"
	"time"
	workflow "workflow-go"
)

func main() {
	client, err := workflow.NewClient("<QSTASH_TOKEN>")
	if err != nil {
		fmt.Printf("failed to create client: %v\n", err)
		return
	}

	// Fetch all logs (with default pagination)
	runs, cursor, err := client.Logs(workflow.LogsOptions{})
	if err != nil {
		fmt.Printf("failed to fetch logs: %v\n", err)
		return
	}
	fmt.Printf("fetched %d workflow runs, next cursor: %s\n", len(runs), cursor)

	// Fetch logs with pagination (limit to 10 runs)
	runs, cursor, err = client.Logs(workflow.LogsOptions{
		Count: 10,
	})
	if err != nil {
		fmt.Printf("failed to fetch logs: %v\n", err)
		return
	}
	fmt.Printf("fetched %d workflow runs, next cursor: %s\n", len(runs), cursor)

	// Fetch the next page of logs using cursor
	if cursor != "" {
		runs, cursor, err = client.Logs(workflow.LogsOptions{
			Cursor: cursor,
			Count:  10,
		})
		if err != nil {
			fmt.Printf("failed to fetch next page of logs: %v\n", err)
			return
		}
		fmt.Printf("fetched next page with %d workflow runs, next cursor: %s\n", len(runs), cursor)
	}

	// Fetch logs for a specific workflow run
	runs, _, err = client.Logs(workflow.LogsOptions{
		Filter: workflow.LogFilter{
			RunId: "workflow-run-id",
		},
	})
	if err != nil {
		fmt.Printf("failed to fetch logs for specific run: %v\n", err)
		return
	}
	if len(runs) > 0 {
		fmt.Printf("found workflow run with ID %s, state: %s\n", runs[0].WorkflowRunId, runs[0].WorkflowState)
	} else {
		fmt.Println("no workflow run found with that ID")
	}

	// Fetch logs for workflows with a specific URL
	runs, _, err = client.Logs(workflow.LogsOptions{
		Filter: workflow.LogFilter{
			Url: "https://example.com/api/workflow",
		},
	})
	if err != nil {
		fmt.Printf("failed to fetch logs for URL: %v\n", err)
		return
	}
	fmt.Printf("found %d workflow runs for the specified URL\n", len(runs))

	// Fetch a specific workflow run by creation time
	// This is because there might be multiple runs with the same ID
	runs, _, err = client.Logs(workflow.LogsOptions{
		Filter: workflow.LogFilter{
			RunId:     "workflow-run-id",
			CreatedAt: time.Now(),
		},
	})
	if err != nil {
		fmt.Printf("failed to fetch logs by creation time: %v\n", err)
		return
	}
	fmt.Printf("found %d workflow runs created after the specified time\n", len(runs))

	// Fetch logs for workflows in a specific state
	runs, _, err = client.Logs(workflow.LogsOptions{
		Filter: workflow.LogFilter{
			State: "RUN_SUCCESS",
		},
	})
	if err != nil {
		fmt.Printf("failed to fetch logs by state: %v\n", err)
		return
	}
	fmt.Printf("found %d completed workflow runs\n", len(runs))

	// Fetch logs with combined filters
	runs, _, err = client.Logs(workflow.LogsOptions{
		Count: 5,
		Filter: workflow.LogFilter{
			State: "RUN_FAILED",
			Url:   "https://example.com/api/workflow",
		},
	})
	if err != nil {
		fmt.Printf("failed to fetch logs with combined filters: %v\n", err)
		return
	}
	fmt.Printf("found %d failed workflow runs for the specified URL\n", len(runs))

	// Print detailed information about the first run (if available)
	if len(runs) > 0 {
		run := runs[0]
		fmt.Println("\nDetailed information about the first workflow run:")
		fmt.Printf("Run ID: %s\n", run.WorkflowRunId)
		fmt.Printf("URL: %s\n", run.WorkflowUrl)
		fmt.Printf("State: %s\n", run.WorkflowState)
		fmt.Printf("Created At: %d\n", run.WorkflowRunCreatedAt)
		fmt.Printf("Completed At: %d\n", run.WorkflowRunCompletedAt)
		fmt.Printf("Number of step groups: %d\n", len(run.GroupedSteps))
	}
}
