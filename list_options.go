package form3

import (
	"net/url"
	"strconv"
)

func getClientWithOptions(endpoint string, opts []ClientOption) *Client {
	cl := &Client{
		baseURL:           endpoint,
		maxRequestTimeout: defaultRequestTimeout,
		validateAccount:   getValidateAccount(),
	}

	for _, opt := range opts {
		opt(cl)
	}
	return cl
}

// ListOption defines an optional parameter type for list calls.
type ListOption func(params url.Values)

// WithPageNo can be used with a list call to pass in the desired page number.
func WithPageNo(page int) ListOption {
	return func(params url.Values) {
		params.Set("page[number]", strconv.Itoa(page))
	}
}

// WithPageSize can be used with a list call to pass in the desired page size.
func WithPageSize(size int) ListOption {
	return func(params url.Values) {
		params.Set("page[size]", strconv.Itoa(size))
	}
}

func applyOptions(params url.Values, options []ListOption) {
	for _, option := range options {
		option(params)
	}
}
