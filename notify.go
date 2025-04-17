package workflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// Notify notifies waiters waiting for the given event ID.
// Event data is expected to be either a string or a JSON serializable object.
// It returns the list of waiters that were notified.
func (c *Client) Notify(eventId string, eventData any) ([]NotifyMessage, error) {
	body := ""
	if reflect.ValueOf(eventData).Kind() != reflect.String {
		data, err := json.Marshal(eventData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall event data: %w", err)
		}
		body = string(data)
	} else {
		body = eventData.(string)
	}
	req := requestOptions{
		method: http.MethodPost,
		path:   []string{"v2", "notify", eventId},
		body:   body,
	}
	resp, _, err := c.do(req)
	if err != nil {
		return nil, err
	}
	r, err := parse[[]NotifyMessage](resp)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type NotifyMessage struct {
	Waiter    Waiter
	MessageId string
	Error     error
}
