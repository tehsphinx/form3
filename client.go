package form3

import (
	"net/http"
	"time"
)

const defaultRequestTimeout = 30 * time.Second

// NewClient creates a new form3 API client. `baseurl` should be scheme, domain and port of the API server.
// Use one or more ClientOption to further configure the client.
func NewClient(endpoint string, opts ...ClientOption) *Client {
	cl := &Client{
		baseURL:           endpoint,
		maxRequestTimeout: defaultRequestTimeout,
		validateAccount:   getValidateAccount(),
	}

	for _, opt := range opts {
		opt(cl)
	}

	cl.client = &http.Client{
		Timeout: cl.maxRequestTimeout,
	}
	return cl
}

// Client implements a client for the form3 API.
type Client struct {
	client *http.Client
	// base url (scheme + domain + port) of the api server
	baseURL string
	// max time limit for all requests
	maxRequestTimeout time.Duration
	// enables debug output
	enableDbg bool
	// validate function for account
	validateAccount func(attr *Account) error
}
