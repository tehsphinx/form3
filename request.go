package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tehsphinx/dbg"
)

// custom errors for http status codes
var (
	ErrNotFound = errors.New("specified resource does not exist")
	ErrConflict = errors.New("specified version incorrect")
)

const (
	typeAccounts attrType = "accounts"
)

type request struct {
	Data requestData `json:"data"`
}

type requestData struct {
	Type           attrType    `json:"type"`
	ID             string      `json:"id"`
	OrganisationID string      `json:"organisation_id"`
	Attributes     interface{} `json:"attributes"`
}

type response struct {
	Data  responseData      `json:"data"`
	Links map[string]string `json:"links"`
}

type listResponse struct {
	Data  []responseData    `json:"data"`
	Links map[string]string `json:"links"`
}

type responseData struct {
	Type           attrType        `json:"type"`
	ID             string          `json:"id"`
	OrganisationID string          `json:"organisation_id"`
	Version        int             `json:"version"`
	CreatedOn      time.Time       `json:"created_on"`
	ModifiedOn     time.Time       `json:"modified_on"`
	Attributes     json.RawMessage `json:"attributes"`
}

/*
request function with the options was created with usability by the API developer in mind.
  It adds some extra heap allocations though.

The function also tries to keep most of the logic in one place. This has pros and cons. One con is its complexity,
  handling many edge cases in one place. A pro: there is only one place to change things.
*/
func (s *Client) request(ctx context.Context, url string, options ...reqOption) error {
	opts := applyReqestOptions(options)

	// marshal request body if request attributes are provided.
	var body []byte
	if opts.reqAttr != nil {
		reqObj := request{
			Data: requestData{
				Type:           opts.attrType,
				ID:             opts.uid,
				OrganisationID: opts.orgID,
				Attributes:     opts.reqAttr,
			},
		}

		var err error
		body, err = json.Marshal(reqObj)
		if err != nil {
			return fmt.Errorf("marshalling request failed: %w", err)
		}
	}

	// build and execute request
	if s.enableDbg {
		dbg.Green(opts.method, url)
		if len(body) != 0 {
			dbg.Blue(string(body))
		}
	}
	req, err := http.NewRequestWithContext(ctx, opts.method, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("invalid request url: %w", err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != opts.statusOK {
		err = errFromStatusCode(resp.StatusCode)
		return fmt.Errorf("unexpected response status %d (%s): %w", resp.StatusCode, resp.Status, err)
	}

	// read and unmarshall response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	if s.enableDbg && len(body) != 0 {
		dbg.Cyan(string(body))
	}

	if opts.response != nil {
		return unmarshalResponse(body, opts)
	}
	if opts.listResponse != nil {
		return unmarshalListResponse(body, opts)
	}
	return nil
}

func unmarshalResponse(body []byte, opts *reqOptions) error {
	if err := json.Unmarshal(body, &opts.response); err != nil {
		return fmt.Errorf("unmarshalling response failed: %w", err)
	}
	if err := json.Unmarshal(opts.response.Data.Attributes, &opts.respAttr); err != nil {
		return fmt.Errorf("unmarshalling response failed: %w", err)
	}
	opts.respAttr.fillFromResponse(opts.response.Data)
	return nil
}

func unmarshalListResponse(body []byte, opts *reqOptions) error {
	if err := json.Unmarshal(body, &opts.listResponse); err != nil {
		return fmt.Errorf("unmarshalling response failed: %w", err)
	}
	for _, item := range opts.listResponse.Data {
		dest := opts.factory()
		if err := json.Unmarshal(item.Attributes, &dest); err != nil {
			return fmt.Errorf("unmarshalling response failed: %w", err)
		}

		dest.fillFromResponse(item)
		opts.callback(dest)
	}
	return nil
}

func errFromStatusCode(statusCode int) error {
	switch statusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusConflict:
		return ErrConflict
	}
	return nil
}
