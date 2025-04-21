package workflow

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Options struct {
	// Url is the address of the QStash.
	// If not provided, the default production URL is used.
	Url string

	// Token is the authentication token for the QStash.
	Token string

	// Client is the HTTP client used to make requests to the QStash server.
	// If not provided, a default HTTP client is used.
	Client *http.Client
}

func (o *Options) validate() error {
	if o.Token == "" {
		return fmt.Errorf("token is empty")
	}
	return nil
}

type Client struct {
	url    string
	token  string
	client *http.Client
}

func newClient(opts Options) (*Client, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	if opts.Url == "" {
		opts.Url = productionUrl
	} else if !strings.HasSuffix(opts.Url, "/") {
		opts.Url += "/"
	}
	if opts.Client == nil {
		opts.Client = http.DefaultClient
	}
	return &Client{
		url:    opts.Url,
		token:  fmt.Sprintf("Bearer %s", opts.Token),
		client: opts.Client,
	}, nil
}

func NewClientWith(opts Options) (*Client, error) {
	return newClient(opts)
}

func NewClient(token string) (*Client, error) {
	return newClient(Options{
		Token: token,
	})
}

func NewClientWithEnv() (*Client, error) {
	return newClient(Options{
		Url:   os.Getenv("QSTASH_URL"),
		Token: os.Getenv("QSTASH_TOKEN"),
	})
}

type requestOptions struct {
	method string
	path   []string
	body   []byte
	header http.Header
	params url.Values
}

func (c *Client) do(opts requestOptions) ([]byte, int, error) {
	request, err := http.NewRequest(opts.method, fmt.Sprintf("%s%s", c.url, strings.Join(opts.path, "/")), bytes.NewReader(opts.body))
	if err != nil {
		return nil, -1, err
	}
	if opts.params != nil {
		request.URL.RawQuery = opts.params.Encode()
	}
	header := opts.header
	if header == nil {
		header = http.Header{}
	}
	header.Set(authorizationHeader, c.token)
	request.Header = header
	res, err := c.client.Do(request)
	if err != nil {
		return nil, -1, err
	}
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, -1, err
	}
	if res.StatusCode >= http.StatusMultipleChoices {
		var rErr restError
		if err = json.Unmarshal(response, &rErr); err != nil {
			return response, res.StatusCode, err
		}
		return response, res.StatusCode, errors.New(rErr.Error)
	}
	return response, res.StatusCode, nil
}

func parse[T any](data []byte) (t T, err error) {
	err = json.Unmarshal(data, &t)
	return
}

type restError struct {
	Error string `json:"error"`
}
