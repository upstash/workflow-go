package workflow

import (
	"net/http"
	"net/url"
	"strconv"
)

type LogFilter struct {
	RunId     string
	Url       string
	CreatedAt int64
	State     string
}

type LogsOptions struct {
	Cursor string
	Count  int
	Filter LogFilter
}

func (l *LogsOptions) params() url.Values {
	params := url.Values{}
	if l.Cursor != "" {
		params.Set("cursor", l.Cursor)
	}
	if l.Count > 0 {
		params.Set("count", strconv.Itoa(l.Count))
	}
	if l.Filter.RunId != "" {
		params.Set("workflowRunId", l.Filter.RunId)
	}
	if l.Filter.State != "" {
		params.Set("state", l.Filter.State)
	}
	if l.Filter.Url != "" {
		params.Set("workflowUrl", l.Filter.Url)
	}
	if l.Filter.CreatedAt != 0 {
		params.Set("workflowCreatedAt", strconv.FormatInt(l.Filter.CreatedAt, 10))
	}
	return params
}

type listRunsResponse struct {
	Cursor string `json:"cursor,omitempty"`
	Runs   []Run  `json:"runs"`
}

// Logs returns a list of workflow runs.
//
// The returned list is sorted by descending creation time.
// If a cursor is returned, it can be used to fetch the next page of results.
//
// The returned runs are a subset of the possible runs that match the filter.
// To get more results, call this function again with the returned cursor.
func (c *Client) Logs(opts LogsOptions) ([]Run, string, error) {
	req := requestOptions{
		method: http.MethodGet,
		path:   []string{"v2", "workflows", "logs"},
		params: opts.params(),
	}
	resp, _, err := c.do(req)
	if err != nil {
		return nil, "", err
	}
	events, err := parse[listRunsResponse](resp)
	if err != nil {
		return nil, "", err
	}
	return events.Runs, events.Cursor, nil
}
