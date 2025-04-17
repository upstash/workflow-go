package workflow_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	workflow "workflow-go"
)

func TestNotify_WithoutWaiter(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)

	messages, err := client.Notify(uuid.NewString(), uuid.NewString())
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
		Body: map[string]any{
			"eventId":           eventId,
			"expectedEventData": expectedEventData,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, runId)

	assert.Eventually(t, func() bool {
		subT := &testing.T{}

		runs, cursor, err := client.Logs(workflow.LogsOptions{
			Filter: workflow.LogFilter{
				RunId: runId,
			},
		})
		assert.NoError(subT, err)
		assert.Empty(subT, cursor)
		if len(runs) != 1 {
			return false
		}
		run := runs[0]
		assert.Len(subT, run.GroupedSteps, 3)
		if len(run.GroupedSteps[len(run.GroupedSteps)-1].Steps) != 1 {
			return false
		}
		lastStep := run.GroupedSteps[len(run.GroupedSteps)-1].Steps[0]
		assert.Equal(subT, lastStep.StepType, "Wait")
		return !subT.Failed()
	}, time.Minute, time.Millisecond*100)

	messages, err := client.Notify(eventId, expectedEventData)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)

	notified := messages[0]
	assert.Equal(t, notified.Waiter.Url, waitForEvent)
	assert.NoError(t, notified.Error)

	waitUntilRunState(t, client, runId, "RUN_SUCCESS")
}
