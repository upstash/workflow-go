package workflow_test

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
	workflow "workflow-go"
)

func TestTrigger(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, err = client.CancelAll()
		assert.NoError(t, err)
	})

	t.Run("basic trigger", func(t *testing.T) {
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url: simpleWorkflowUrl,
			Body: map[string]string{
				"name": "John Doe",
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, runId)

		waitUntilRunState(t, client, runId, "RUN_SUCCESS")
	})

	t.Run("custom run id", func(t *testing.T) {
		customId := uuid.NewString()
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url:   simpleWorkflowUrl,
			RunId: customId,
			Body: map[string]string{
				"name": "Custom ID Test",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, customId, runId)

		waitUntilRunState(t, client, runId, "RUN_SUCCESS")
	})

	t.Run("with retries", func(t *testing.T) {
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url:     failingWorkflow,
			Retries: workflow.Retry(0),
			Body: map[string]string{
				"name": "Retry Test",
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, runId)

		waitUntilRunState(t, client, runId, "RUN_FAILED")
	})

	// Test with custom headers
	t.Run("with custom headers", func(t *testing.T) {
		customHeaders := http.Header{}
		customHeaders.Set("X-Custom-Header", "test-value")
		customHeaders.Set("X-Another-Header", "another-value")

		runId, err := client.Trigger(workflow.TriggerOptions{
			Url:    simpleWorkflowUrl,
			Header: customHeaders,
			Body: map[string]string{
				"name": "Custom Headers Test",
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, runId)

		assert.Eventually(t, func() bool {
			subT := &testing.T{}

			runs, _, err := client.Logs(workflow.LogsOptions{Filter: workflow.LogFilter{RunId: runId}})
			assert.NoError(subT, err)

			if assert.Len(subT, runs, 1) {
				run := runs[0]
				if assert.GreaterOrEqual(subT, len(run.GroupedSteps), 1) {
					step := run.GroupedSteps[0]

					assert.Equal(subT, step.Steps[0].Headers.Get("X-Custom-Header"), "test-value")
					assert.Equal(subT, step.Steps[0].Headers.Get("X-Another-Header"), "another-value")
				}
			}
			return !subT.Failed()
		}, time.Minute, time.Millisecond*100)
	})

	t.Run("string body", func(t *testing.T) {
		body := "plain text body content"
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url:  simpleWorkflowUrl,
			Body: "plain text body content",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, runId)

		assertRequestPayload(t, client, runId, body)
		waitUntilRunState(t, client, runId, "RUN_SUCCESS")
	})

	t.Run("struct body", func(t *testing.T) {
		type TestStruct struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}
		requestPayload := TestStruct{
			Name:  "Struct Test",
			Value: 42,
		}
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url:  simpleWorkflowUrl,
			Body: requestPayload,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, runId)

		data, err := json.Marshal(requestPayload)
		assert.NoError(t, err)

		assertRequestPayload(t, client, runId, string(data))
		waitUntilRunState(t, client, runId, "RUN_SUCCESS")
	})

	t.Run("nil body", func(t *testing.T) {
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url:  simpleWorkflowUrl,
			Body: nil,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, runId)

		assertRequestPayload(t, client, runId, "")
		waitUntilRunState(t, client, runId, "RUN_SUCCESS")
	})

	t.Run("empty url", func(t *testing.T) {
		_, err := client.Trigger(workflow.TriggerOptions{
			Body: map[string]string{
				"name": "Error Test",
			},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "url is required")
	})
}
