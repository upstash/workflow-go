package workflow_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/upstash/workflow-go"
	"testing"
)

func TestCancel(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)

	runId, err := client.Trigger(workflow.TriggerOptions{
		Url: sleepTenDaysWorkflow,
	})
	assert.NoError(t, err)

	err = client.Cancel(runId)
	assert.NoError(t, err)
}

func TestCancel_UnknownRun(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)

	err = client.Cancel("unknown-run-id")
	assert.ErrorContains(t, err, "failed to cancel workflow run")
}

func TestClient_CancelMany(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, err = client.CancelAll()
		assert.NoError(t, err)
	})

	count := 5
	runIds := make([]string, count)
	for i := 0; i < count; i++ {
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url: sleepTenDaysWorkflow,
		})
		assert.NoError(t, err)
		runIds[i] = runId
	}

	canceled, err := client.CancelMany(runIds[:count-1])
	assert.NoError(t, err)
	assert.Equal(t, canceled, count-1)
}

func TestClient_CancelAll(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)

	count := 5
	runIds := make([]string, count)
	for i := 0; i < count; i++ {
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url: sleepTenDaysWorkflow,
		})
		assert.NoError(t, err)
		runIds[i] = runId
	}

	canceled, err := client.CancelAll()
	assert.NoError(t, err)
	assert.Equal(t, canceled, count)
}

func TestClient_CancelWithFilter(t *testing.T) {
	client, err := workflow.NewClientWithEnv()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, err = client.CancelAll()
		assert.NoError(t, err)
	})

	count := 5
	typeARunIds := make([]string, count)
	for i := 0; i < count; i++ {
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url: sleepTenDaysWorkflow,
		})
		assert.NoError(t, err)
		typeARunIds[i] = runId
	}

	typeBRunIds := make([]string, count)
	for i := 0; i < count; i++ {
		runId, err := client.Trigger(workflow.TriggerOptions{
			Url: sleepOneDayWorkflow,
		})
		assert.NoError(t, err)
		typeBRunIds[i] = runId
	}

	canceled, err := client.CancelWithFilter(workflow.CancelFilter{
		UrlStartingWith: sleepOneDayWorkflow,
	})
	assert.NoError(t, err)
	assert.Equal(t, canceled, count)

	canceled, err = client.CancelWithFilter(workflow.CancelFilter{
		RunIds:          typeBRunIds,
		UrlStartingWith: sleepTenDaysWorkflow,
	})
	assert.NoError(t, err)
	assert.Zero(t, canceled)

	canceled, err = client.CancelWithFilter(workflow.CancelFilter{
		UrlStartingWith: sleepTenDaysWorkflow,
	})
	assert.NoError(t, err)
	assert.Equal(t, canceled, count)
}
