package workflow

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CancelFilter struct {
	UrlStartingWith string
	RunIds          []string
}
type cancelRequest struct {
	RunIds []string `json:"workflowRunIds,omitempty"`
	Url    string   `json:"workflowUrl,omitempty"`
}

type cancelResponse struct {
	Cancelled int `json:"cancelled"`
}

// Cancel cancels an ongoing workflow run.
// If a non-ongoing run is passed, it will fail.
func (c *Client) Cancel(runId string) error {
	canceled, err := c.cancel(cancelRequest{RunIds: []string{runId}})
	if err != nil {
		return err
	}
	if canceled != 1 {
		return fmt.Errorf("failed to cancel workflow run %s", runId)
	}
	return nil
}

// CancelMany cancels the given workflow runs.
// It returns the total number of canceled workflow runs.
func (c *Client) CancelMany(runIds []string) (int, error) {
	return c.cancel(cancelRequest{RunIds: runIds})
}

// CancelAll cancels all the ongoing workflow runs.
func (c *Client) CancelAll() (int, error) {
	return c.cancel(cancelRequest{})
}

// CancelWithFilter cancels workflow runs that match the specified filter criteria.
// It cancels runs whose URLs start with the given prefix and/or whose run IDs are specified.
// Returns the number of canceled workflow runs or an error if the operation fails.
func (c *Client) CancelWithFilter(filter CancelFilter) (int, error) {
	return c.cancel(cancelRequest{Url: filter.UrlStartingWith, RunIds: filter.RunIds})
}

func (c *Client) cancel(req cancelRequest) (int, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshall request: %w", err)
	}
	resp, _, err := c.do(requestOptions{
		method: http.MethodDelete,
		path:   []string{"v2", "workflows", "runs"},
		body:   string(data),
		header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	})
	if err != nil {
		return 0, err
	}
	canceled, err := parse[cancelResponse](resp)
	if err != nil {
		return 0, err
	}
	return canceled.Cancelled, err
}
