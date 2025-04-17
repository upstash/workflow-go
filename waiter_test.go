package workflow_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	workflow "workflow-go"
)

func TestWaiters(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, err = client.CancelAll()
		assert.NoError(t, err)
	})

	eventId := uuid.NewString()
	runId, err := client.Trigger(workflow.TriggerOptions{
		Url: waitForEvent,
		Body: map[string]any{
			"eventId": eventId,
			"expectedEventData": map[string]string{
				"test": "data",
			},
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, runId)

	waitUntilWaitStep(t, client, runId)

	waiters, err := client.Waiters(eventId)
	assert.NoError(t, err)
	assert.Len(t, waiters, 1)

	for _, waiter := range waiters {
		assert.Equal(t, waiter.Url, waitForEvent)
		assert.Equal(t, waiter.Headers.Get("Upstash-Workflow-Runid"), runId)
	}

	resp, err := client.Notify(eventId, map[string]string{
		"test": "data",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	waitUntilRunState(t, client, runId, "RUN_SUCCESS")

	waiters, err = client.Waiters(eventId)
	assert.NoError(t, err)
	assert.Empty(t, waiters)
}

func TestWaiters_NoWaiters(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)

	eventId := uuid.NewString()

	waiters, err := client.Waiters(eventId)
	assert.NoError(t, err)
	assert.Empty(t, waiters)
}

func TestWaiters_MultipleEventIdsAndRuns(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, err = client.CancelAll()
		assert.NoError(t, err)
	})

	eventIds := []string{}
	runIds := []string{}
	for i := 0; i < 5; i++ {
		eventId := uuid.NewString()

		for j := 0; j < 2; j++ {
			runId, err := client.Trigger(workflow.TriggerOptions{
				Url: waitForEvent,
				Body: map[string]any{
					"eventId": eventId,
					"expectedEventData": map[string]string{
						"test": eventId,
					},
				},
			})
			assert.NoError(t, err)
			assert.NotEmpty(t, runId)
			runIds = append(runIds, runId)
		}
		eventIds = append(eventIds, eventId)
	}

	for _, runId := range runIds {
		waitUntilWaitStep(t, client, runId)
	}

	for _, eventId := range eventIds {
		waiters, err := client.Waiters(eventId)
		assert.NoError(t, err)
		assert.Len(t, waiters, 2)
	}
}
