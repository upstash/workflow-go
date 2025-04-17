package workflow

import "net/http"

type RunResponse struct {
	Body       string `json:"body,omitempty"`
	IsCanceled bool   `json:"isCanceled"`
	IsFailed   bool   `json:"isFailed"`
}

type InvokerContext struct {
	WorkflowRunId     string `json:"workflowRunId,omitempty"`
	WorkflowUrl       string `json:"workflowUrl,omitempty"`
	WorkflowCreatedAt int64  `json:"workflowRunCreatedAt,omitempty"`
}

type FailureFunction struct {
	MessageId    string      `json:"messageId"`
	Url          string      `json:"url"`
	State        string      `json:"state"`
	FailHeaders  http.Header `json:"failHeaders"`
	FailStatus   int         `json:"failStatus"`
	FailResponse string      `json:"failResponse"`
	DlqId        string      `json:"dlqId"`
}

type StepInfo struct {
	StepId     int64       `json:"stepId,omitempty"`
	StepName   string      `json:"stepName,omitempty"`
	StepType   string      `json:"stepType,omitempty"`
	CallType   string      `json:"callType,omitempty"`
	MessageId  string      `json:"messageId"`
	Headers    http.Header `json:"headers,omitempty"`
	TargetStep int64       `json:"targetStep,omitempty"`
	Out        string      `json:"out,omitempty"`
	Concurrent int64       `json:"concurrent,omitempty"`
	State      string      `json:"state,omitempty"`
	CreatedAt  int64       `json:"createdAt,omitempty"`

	SleepUntil int64 `json:"sleepUntil,omitempty"`
	SleepFor   int64 `json:"sleepFor,omitempty"`

	CallUrl             string      `json:"callUrl,omitempty"`
	CallMethod          string      `json:"callMethod,omitempty"`
	CallBody            string      `json:"callBody,omitempty"`
	CallHeaders         http.Header `json:"callHeaders,omitempty"`
	CallResponseStatus  int64       `json:"callResponseStatus,omitempty"`
	CallResponseBody    string      `json:"callResponseBody,omitempty"`
	CallResponseHeaders http.Header `json:"callResponseHeaders,omitempty"`

	WaitEventId         string `json:"waitEventId,omitempty"`
	WaitTimeoutDuration string `json:"waitTimeoutDuration,omitempty"`
	WaitTimeoutDeadline int64  `json:"waitTimeoutDeadline,omitempty"`
	WaitTimeout         bool   `json:"waitTimeout,omitempty"`

	InvokedWorkflowRunId      string      `json:"invokedWorkflowRunId,omitempty"`
	InvokedWorkflowUrl        string      `json:"invokedWorkflowUrl,omitempty"`
	InvokedWorkflowCreatedAt  int64       `json:"invokedWorkflowCreatedAt,omitempty"`
	InvokedWorkflowRunBody    string      `json:"invokedWorkflowRunBody,omitempty"`
	InvokedWorkflowRunHeaders http.Header `json:"invokedWorkflowRunHeaders,omitempty"`
}

type GroupedSteps struct {
	Steps []StepInfo `json:"steps"`
	Type  string     `json:"type"`
}

type Run struct {
	WorkflowRunId          string `json:"workflowRunId"`
	WorkflowUrl            string `json:"workflowUrl"`
	WorkflowState          string `json:"workflowState"`
	WorkflowRunCreatedAt   int64  `json:"workflowRunCreatedAt"`
	WorkflowRunCompletedAt int64  `json:"workflowRunCompletedAt,omitempty"`

	GroupedSteps        []GroupedSteps   `json:"steps,omitempty"`
	WorkflowRunResponse string           `json:"workflowRunResponse,omitempty"`
	Invoker             *InvokerContext  `json:"invoker,omitempty"`
	FailureFunction     *FailureFunction `json:"failureFunction,omitempty"`
}

type Waiter struct {
	Url            string      `json:"url"`
	Headers        http.Header `json:"headers,omitempty"`
	Deadline       int64       `json:"deadline"`
	TimeoutBody    []byte      `json:"timeoutBody,omitempty"`
	TimeoutUrl     string      `json:"timeoutUrl,omitempty"`
	TimeoutHeaders http.Header `json:"timeoutHeaders,omitempty"`
}
