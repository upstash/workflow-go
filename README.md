# Upstash Workflow Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/upstash/workflow-go.svg)](https://pkg.go.dev/github.com/upstash/workflow-go)

> [!NOTE]  
> **This project is in GA Stage.**
>
> The Upstash Professional Support fully covers this project. It receives regular updates, and bug fixes.
> The Upstash team is committed to maintaining and improving its functionality.

**Upstash Workflow** lets you write durable, reliable and performant serverless functions. Get delivery guarantees, automatic retries on failure, scheduling and more without managing any infrastructure.

This is the HTTP-based Go client for [Upstash](https://upstash.com/) Workflow.

Note that this SDK only offers client-side workflow functions for managing workflows; it does not include server-side functionality for writing workflows in Golang.

## Documentation

- [**Reference Documentation**](https://upstash.com/docs/workflow/overall/getstarted)
- [**Quickstart**](https://upstash.com/docs/workflow/quickstarts)
- [**API Reference**](https://pkg.go.dev/github.com/upstash/workflow-go)


## Installation

Use `go get` to install the Upstash Workflow package:
```bash
go get github.com/upstash/workflow-go
```

Import the Upstash Workflow package in your project:

```go
import "github.com/upstash/workflow-go"
```


## Usage

The `QSTASH_TOKEN` is required to initialize an Upstash Workflow client. 
Find your credentials in the console dashboard at [Upstash Console](https://console.upstash.com/qstash).

```go
import (
	"github.com/upstash/workflow-go"
)

func main() {
	client := workflow.NewClient("<QSTASH_TOKEN>")
}
```

Alternatively, you can set the following environment variables:

```shell
QSTASH_URL="<QSTASH_URL>"
QSTASH_TOKEN="<QSTASH_TOKEN>"
```

and then create the client by using:

```go
import (
	"github.com/upstash/workflow-go"
)

func main() {
	client := workflow.NewClientWithEnv()
}
```
#### Using a custom HTTP client

By default, `http.DefaultClient` will be used for doing requests. It is possible to use a custom HTTP client by passing it in the options while constructing the client.

```go
import (
	"net/http"

	"github.com/upstash/workflow-go"
)

func main() {
	opts := workflow.Options{
		Token:  "<QSTASH_TOKEN>",
		Client: &http.Client{},
	}
	client := workflow.NewClientWith(opts)
}
```

### Trigger a workflow run

Start a new workflow run with provided options.

```go
runID, err := client.Trigger(workflow.TriggerOptions{
    Url: "https://your-workflow-endpoint.com/api/process"
    Body: map[string]interface{}{
        "user_id":    12345,
        "action":     "purchase",
        "items":      []string{"item1", "item2", "item3"},
        "total_cost": 99.99,
    },
})
if err != nil {
    // handle err
}
```

### Notify Events

Send a notify message to workflows waiting for a specific event.

```go

messages, err := client.Notify("event-id", map[string]string{
		"userId": "testUser",
})
if err != nil {
    // handle err
}
```

### Cancel Workflows

Cancel one or more ongoing workflow runs.

```go
err := client.Cancel("workflow-run-id")
if err != nil {
    // handle err
}

canceled, err := client.CancelMany([]string{"run-a", "run-b", "run-c"})
if err != nil {
    // handle err
}

canceled, err = client.CancelAll()
if err != nil {
    // handle err
}
```

### Fetch waiters

Get the list of workflows waiting on a specific event.

```go
waiters, err := client.Waiters("my-event-id")
if err != nil {
    // handle err
}
```

### Fetch Logs

Get the logs for workflow runs with filtering.

```go
runs, cursor, err := client.Logs(workflow.LogsOptions{})
if err != nil {
    // handle err
}

runs, cursor, err = client.Logs(workflow.LogsOptions{
    Filter: workflow.LogFilter{
        RunId: "workflow-run-id",
    },
})
if err != nil {
	// handle err
}

runs, cursor, err = client.Logs(workflow.LogsOptions{
    Filter: workflow.LogFilter{
        State: "RUN_SUCCESS",
    },
})
if err != nil {
	// handle err
}
```
