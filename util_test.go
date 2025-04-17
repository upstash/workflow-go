package workflow_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
	workflow "workflow-go"
)

var (
	testWorkflowBaseUrl  = os.Getenv("APPLICATION_URL")
	simpleWorkflowUrl    = strings.Join([]string{testWorkflowBaseUrl, "workflows", "simpleWorkflow"}, "/")
	longRunningWorkflow  = strings.Join([]string{testWorkflowBaseUrl, "workflows", "longRunning"}, "/")
	sleepTenDaysWorkflow = strings.Join([]string{testWorkflowBaseUrl, "workflows", "sleepTenDaysWorkflow"}, "/")
	sleepOneDayWorkflow  = strings.Join([]string{testWorkflowBaseUrl, "workflows", "sleepOneDayWorkflow"}, "/")
	waitForEvent         = strings.Join([]string{testWorkflowBaseUrl, "workflows", "waitForEvent"}, "/")
	failingWorkflow      = strings.Join([]string{testWorkflowBaseUrl, "workflows", "failingWorkflow"}, "/")
)

func waitUntilRunState(t *testing.T, client *workflow.Client, runId string, state string) {
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
		assert.Equal(subT, run.WorkflowState, state)
		return !subT.Failed()
	}, time.Minute, time.Millisecond*100)
}

func assertRequestPayload(t *testing.T, client *workflow.Client, runId string, expected string) {
	assert.Eventually(t, func() bool {
		subT := &testing.T{}

		runs, _, err := client.Logs(workflow.LogsOptions{Filter: workflow.LogFilter{RunId: runId}})
		assert.NoError(subT, err)

		if assert.Len(subT, runs, 1) {
			run := runs[0]
			if assert.GreaterOrEqual(subT, len(run.GroupedSteps), 1) {
				step := run.GroupedSteps[0]

				assert.Equal(subT, step.Steps[0].Out, expected)
			}
		}
		return !subT.Failed()
	}, time.Minute, time.Millisecond*100)
}

func waitUntilWaitStep(t *testing.T, client *workflow.Client, runId string) {
	assert.Eventually(t, func() bool {
		subT := &testing.T{}

		runs, _, err := client.Logs(workflow.LogsOptions{
			Filter: workflow.LogFilter{
				RunId: runId,
			},
		})
		assert.NoError(subT, err)
		if len(runs) != 1 {
			return false
		}
		run := runs[0]
		if len(run.GroupedSteps) < 3 {
			return false
		}
		if len(run.GroupedSteps[len(run.GroupedSteps)-1].Steps) != 1 {
			return false
		}
		lastStep := run.GroupedSteps[len(run.GroupedSteps)-1].Steps[0]
		assert.Equal(subT, lastStep.StepType, "Wait")
		return !subT.Failed()
	}, time.Minute, time.Millisecond*100)
}
