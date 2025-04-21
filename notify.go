package workflow

import (
	"net/http"
)

// Notify notifies waiters waiting for the given event ID.
// Event data is sent to the waiters as is.
// It returns the list of waiters that were notified.
func (c *Client) Notify(eventId string, eventData []byte) ([]NotifyMessage, error) {
	req := requestOptions{
		method: http.MethodPost,
		path:   []string{"v2", "notify", eventId},
		body:   eventData,
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
