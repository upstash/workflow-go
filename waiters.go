package workflow

import "net/http"

// Waiters returns all waiters waiting for the given event ID.
// A waiter is registered when a workflow executes a context.waitForEvent() step.
// A waiter is removed once the user notifies or the timeout duration has expired.
func (c *Client) Waiters(eventId string) ([]Waiter, error) {
	req := requestOptions{
		method: http.MethodGet,
		path:   []string{"v2", "waiters", eventId},
	}
	resp, _, err := c.do(req)
	if err != nil {
		return nil, err
	}
	return parse[[]Waiter](resp)
}
