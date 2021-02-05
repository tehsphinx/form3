package form3

import "net/http"

type reqOption interface {
	apply(opts *reqOptions)
}

type reqOptionFunc func(opts *reqOptions)

func (s reqOptionFunc) apply(opts *reqOptions) {
	s(opts)
}

type reqOptions struct {
	method       string
	orgID        string
	uid          string
	reqAttr      interface{}
	response     *response
	listResponse *listResponse
	respAttr     responseFiller
	factory      func() responseFiller
	callback     func(responseFiller)
	statusOK     int
	attrType     attrType
}

type attrType string

func (s attrType) apply(opts *reqOptions) {
	opts.attrType = s
}

func applyReqestOptions(options []reqOption) *reqOptions {
	opts := &reqOptions{
		method:   http.MethodGet,
		statusOK: http.StatusOK,
	}
	for _, o := range options {
		o.apply(opts)
	}
	return opts
}

// withMethod adds the http method to a request. Defaults to GET.
func withMethod(httpMethod string) reqOptionFunc {
	return func(opts *reqOptions) {
		opts.method = httpMethod
	}
}

// withOrgID adds an organisation ID to the request.
func withOrgID(orgID string) reqOptionFunc {
	return func(opts *reqOptions) {
		opts.orgID = orgID
	}
}

// withUID adds a uid to the request.
func withUID(uid string) reqOptionFunc {
	return func(opts *reqOptions) {
		opts.uid = uid
	}
}

// withReq adds the attributes part to the request body.
func withReq(attributes interface{}) reqOptionFunc {
	return func(opts *reqOptions) {
		opts.reqAttr = attributes
	}
}

// withReq adds the attributes part to the request body.
func withStatusOk(httpStatus int) reqOptionFunc {
	return func(opts *reqOptions) {
		opts.statusOK = httpStatus
	}
}

type responseFiller interface {
	fillFromResponse(resp responseData)
}

// withResp adds a pointer to be filled with the attributes part of the response body.
func withResp(attributes responseFiller) reqOptionFunc {
	return func(opts *reqOptions) {
		opts.respAttr = attributes
		opts.response = &response{
			Data: responseData{},
		}
	}
}

// withListResp adds a pointer to be filled with the attributes part of the response body.
func withListResp(factory func() responseFiller, cb func(responseFiller)) reqOptionFunc {
	return func(opts *reqOptions) {
		opts.listResponse = &listResponse{
			Data: []responseData{},
		}
		opts.factory = factory
		opts.callback = cb
	}
}
