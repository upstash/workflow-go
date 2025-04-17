package workflow_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	workflow "workflow-go"
)

func TestLogs(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)

	var runIds []string
	for i := 0; i < 3; i++ {
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url: simpleWorkflowUrl,
			Body: map[string]string{
				"name": "Run Test User",
				"id":   string(rune(65 + i)),
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, runId)
		runIds = append(runIds, runId)
	}

	// Trigger a long-running workflow for state filtering test
	longRunId, err := client.Trigger(workflow.TriggerOptions{
		Url: longRunningWorkflow,
		Body: map[string]string{
			"name": "Long Running Test",
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, longRunId)
	runIds = append(runIds, longRunId)

	for i := 0; i < 3; i++ {
		waitUntilRunState(t, client, runIds[i], "RUN_SUCCESS")
	}

	t.Run("list all runs", func(t *testing.T) {
		runs, _, err := client.Logs(workflow.LogsOptions{})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(runs), len(runIds))

		actualRunIds := []string{}
		for _, run := range runs {
			actualRunIds = append(actualRunIds, run.WorkflowRunId)
		}
		for _, runId := range runIds {
			assert.Contains(t, actualRunIds, runId)
		}
	})

	t.Run("limit results with count", func(t *testing.T) {
		count := 2
		runs, cursor, err := client.Logs(workflow.LogsOptions{
			Count: count,
		})
		assert.NoError(t, err)
		assert.Len(t, runs, count)
		assert.NotEmpty(t, cursor)
	})

	t.Run("filter by run id", func(t *testing.T) {
		runs, cursor, err := client.Logs(workflow.LogsOptions{
			Filter: workflow.LogFilter{
				RunId: runIds[0],
			},
		})
		assert.NoError(t, err)
		assert.Empty(t, cursor)

		if assert.Len(t, runs, 1) {
			run := runs[0]
			assert.Equal(t, runIds[0], run.WorkflowRunId)
			assert.Equal(t, simpleWorkflowUrl, run.WorkflowUrl)
		}
	})

	t.Run("filter by url", func(t *testing.T) {
		runs, cursor, err := client.Logs(workflow.LogsOptions{
			Filter: workflow.LogFilter{
				Url: simpleWorkflowUrl,
			},
		})
		assert.Empty(t, cursor)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(runs), 3)

		for _, run := range runs {
			assert.Equal(t, simpleWorkflowUrl, run.WorkflowUrl)
		}

		foundCount := 0
		for _, run := range runs {
			for i := 0; i < 3; i++ {
				if run.WorkflowRunId == runIds[i] {
					foundCount++
					break
				}
			}
		}
		assert.GreaterOrEqual(t, foundCount, 3, "Should find all our simple workflow runs")
	})

	t.Run("filter by state", func(t *testing.T) {
		runs, cursor, err := client.Logs(workflow.LogsOptions{
			Filter: workflow.LogFilter{
				State: "RUN_SUCCESS",
			},
		})
		assert.Empty(t, cursor)
		assert.NoError(t, err)
		assert.NotEmpty(t, runs)

		for _, run := range runs {
			assert.Equal(t, "RUN_SUCCESS", run.WorkflowState)
		}

		foundCount := 0
		for _, run := range runs {
			for i := 0; i < 3; i++ {
				if run.WorkflowRunId == runIds[i] {
					foundCount++
					break
				}
			}
		}
		assert.GreaterOrEqual(t, foundCount, 3, "Should find all our completed workflow runs")
	})

	t.Run("pagination with cursor", func(t *testing.T) {
		// First page with 2 items
		firstPageRuns, cursor, err := client.Logs(workflow.LogsOptions{
			Count: 2,
		})
		assert.NoError(t, err)
		assert.Len(t, firstPageRuns, 2)
		assert.NotEmpty(t, cursor)

		// Second page
		secondPageRuns, _, err := client.Logs(workflow.LogsOptions{
			Cursor: cursor,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, secondPageRuns)

		// Verify pages contain different runs
		assert.NotEqual(t, firstPageRuns[0].WorkflowRunId, secondPageRuns[0].WorkflowRunId)
		assert.NotEqual(t, firstPageRuns[1].WorkflowRunId, secondPageRuns[0].WorkflowRunId)
	})

	// Test with created at filter
	t.Run("FilterByCreatedAt", func(t *testing.T) {
		// Use timestamp from 10 seconds before test started
		createdAfter := time.Now().Add(-time.Minute).Unix()

		runs, _, err := client.Logs(workflow.LogsOptions{
			Filter: workflow.LogFilter{
				CreatedAt: createdAfter,
			},
		})
		assert.NoError(t, err)

		// All our runs should be included since we just created them
		foundCount := 0
		for _, run := range runs {
			for _, id := range runIds {
				if run.WorkflowRunId == id {
					foundCount++
					break
				}
			}
			// Verify timestamp
			assert.GreaterOrEqual(t, run.WorkflowRunCreatedAt, createdAfter)
		}
		assert.GreaterOrEqual(t, foundCount, len(runIds), "All our test runs should be found")
	})

	// Test multiple filters combined
	t.Run("CombinedFilters", func(t *testing.T) {
		runs, _, err := client.Logs(workflow.LogsOptions{
			Count: 10,
			Filter: workflow.LogFilter{
				Url:   simpleWorkflowUrl,
				State: "RUN_SUCCESS",
			},
		})
		assert.NoError(t, err)

		// All returned runs should match both filters
		for _, run := range runs {
			assert.Equal(t, simpleWorkflowUrl, run.WorkflowUrl)
			assert.Equal(t, "RUN_SUCCESS", run.WorkflowState)
		}

		// Our completed simple workflows should be in the results
		foundCount := 0
		for _, run := range runs {
			for i := 0; i < 3; i++ {
				if run.WorkflowRunId == runIds[i] {
					foundCount++
					break
				}
			}
		}
		assert.GreaterOrEqual(t, foundCount, 3, "Should find all our completed simple workflow runs")
	})
}
