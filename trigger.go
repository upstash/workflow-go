package workflow

import (
	"fmt"
	"net/http"
	"strconv"
)

type TriggerOptions struct {
	// Url of the workflow
	Url string
	// Body is the request payload of the new workflow run.
	// It is expected to be either a string or JSON serializable object.
	Body any
	// RunId is the id of the new workflow run. If not provided, a random id will be generated
	RunId string
	// Retries is the number of retries if a step fails.
	Retries        *int
	FlowControlKey string
	Rate           *int
	Parallelism    *int
	// Header is the custom headers that will be forwarded to the workflow.
	Header http.Header
}

func (o *TriggerOptions) validate() error {
	if o.Url == "" {
		return fmt.Errorf("url is required")
	}
	if o.RunId == "" {
		o.RunId = newRunId()
	}
	return nil
}

func (o *TriggerOptions) header() http.Header {
	header := http.Header{}
	header.Set(initHeader, "true")
	header.Set(runIdHeader, o.RunId)
	header.Set(urlHeader, o.Url)
	header.Add(featureSetHeader, featureLazyFetch)
	header.Add(featureSetHeader, featureInitialBody)
	if o.FlowControlKey != "" {
		value := ""
		if o.Rate != nil {
			value += "rate=" + strconv.Itoa(*o.Rate)
		}
		if len(value) != 0 {
			value += ","
		}
		if o.Parallelism != nil {
			value += "parallelism=" + strconv.Itoa(*o.Parallelism)
		}
		if len(value) > 0 {
			header.Set("Upstash-Flow-Control-Key", o.FlowControlKey)
			header.Set("Upstash-Flow-Control-Value", value)
		}
	}
	header.Set(fmt.Sprintf("%s%s", forwardPrefix, sdkVersionHeader), sdkVersion)
	for k, v := range o.Header {
		for _, vv := range v {
			header.Add(fmt.Sprintf("%s%s", forwardPrefix, k), vv)
		}
	}
	if o.Retries != nil {
		header.Set("Upstash-Retries", strconv.Itoa(*o.Retries))
	}
	return header
}

type publishResponse struct {
	MessageId string `json:"messageId"`
}

// Trigger starts a new workflow run and returns the workflow run id.
func (c *Client) Trigger(opts TriggerOptions) (runId string, err error) {
	if err = opts.validate(); err != nil {
		return "", fmt.Errorf("failed to validate options: %w", err)
	}
	body, isJSON, err := serializeToStr(opts.Body)
	if err != nil {
		return "", fmt.Errorf("failed to serialize body: %w", err)
	}
	header := opts.header()
	if isJSON {
		header.Set(contentTypeHeader, "application/json")
	}
	req := requestOptions{
		method: http.MethodPost,
		path:   []string{"v2", "publish", opts.Url},
		body:   body,
		header: header,
	}
	resp, _, err := c.do(req)
	if err != nil {
		return "", err
	}
	_, err = parse[publishResponse](resp)
	if err != nil {
		return "", fmt.Errorf("unexpected response: %w", err)
	}
	return opts.RunId, nil
}
