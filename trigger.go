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
	Body []byte
	// RunId is the id of the new workflow run. If not provided, a random id will be generated
	RunId string
	// Retries is the number of retries if a step fails.
	Retries *int
	// FlowControlKey is the key used to control the flow of new steps.
	FlowControlKey string
	// Rate is the number of new starting steps per period.
	Rate int
	// Parallelism defines the maximum number of active steps associated with this FlowControlKey.
	Parallelism int
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
		givenRate := o.Rate != 0
		if givenRate {
			value += fmt.Sprintf("rate=%d", o.Rate)
		}
		if o.Parallelism != 0 {
			if givenRate {
				value += ","
			}
			value += fmt.Sprintf("parallelism=%d", o.Parallelism)
		}
		if len(value) > 0 {
			header.Set(flowControlKeyHeader, o.FlowControlKey)
			header.Set(flowControlValueHeader, value)
		}
	}
	if contentType := o.Header.Get(contentTypeHeader); contentType != "" {
		header.Set(contentTypeHeader, contentType)
		o.Header.Del(contentTypeHeader)
	}
	header.Set(fmt.Sprintf("%s%s", forwardPrefix, sdkVersionHeader), sdkVersion)
	for k, v := range o.Header {
		for _, vv := range v {
			header.Add(fmt.Sprintf("%s%s", forwardPrefix, k), vv)
		}
	}
	if o.Retries != nil {
		header.Set(retriesHeader, strconv.Itoa(*o.Retries))
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
	header := opts.header()
	req := requestOptions{
		method: http.MethodPost,
		path:   []string{"v2", "publish", opts.Url},
		body:   opts.Body,
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
