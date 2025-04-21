package workflow_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	workflow "workflow-go"
)

func TestNotify_WithoutWaiter(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)

	messages, err := client.Notify(uuid.NewString(), []byte(uuid.NewString()))
	assert.NoError(t, err)
	assert.Len(t, messages, 0)
}

func TestNotify(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, err = client.CancelAll()
		assert.NoError(t, err)
	})

	eventId := uuid.NewString()
	expectedEventData := map[string]string{
		"uuid": uuid.NewString(),
	}
	runId, err := client.Trigger(workflow.TriggerOptions{
		Url: waitForEvent,
		Body: jsonMarshall(t, map[string]any{
			"eventId":           eventId,
			"expectedEventData": expectedEventData,
		}),
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, runId)

	waitUntilWaitStep(t, client, runId)

	messages, err := client.Notify(eventId, jsonMarshall(t, expectedEventData))
	assert.NoError(t, err)
	assert.Len(t, messages, 1)

	notified := messages[0]
	assert.Equal(t, notified.Waiter.Url, waitForEvent)
	assert.NoError(t, notified.Error)

	waitUntilRunState(t, client, runId, "RUN_SUCCESS")
}
