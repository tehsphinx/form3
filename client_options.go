package form3

import "time"

// ClientOption defines an optional parameter for creating a form3.NewClient client.
type ClientOption func(cl *Client)

// WithRequestTimeout sets a maximum request timeout on the client for all requests.
// Individual requests can be timeout out by context in a shorter time. The timeout
// set here can't be extended with individual timeouts.
func WithRequestTimeout(timeout time.Duration) ClientOption {
	return func(cl *Client) {
		cl.maxRequestTimeout = timeout
	}
}

// WithDebug enables colorful debugging of the communication.
func WithDebug() ClientOption {
	return func(cl *Client) {
		cl.enableDbg = true
	}
}
