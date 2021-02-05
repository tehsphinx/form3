package form3

import (
	"time"
)

// installs the response data in a data type without exposing the
// data in json requests and also not making it settable.
type baseAttr struct {
	data responseData
}

// Type returns the type of the data record.
func (s baseAttr) Type() string {
	return string(s.data.Type)
}

// ID returns the uuid of the data record.
func (s baseAttr) ID() string {
	return s.data.ID
}

// OrganisationID returns the organisation id the data record belongs to.
func (s baseAttr) OrganisationID() string {
	return s.data.OrganisationID
}

// Version returns the version of the data record.
func (s baseAttr) Version() int {
	return s.data.Version
}

// CreatedOn returns the creation date of the data record.
func (s baseAttr) CreatedOn() time.Time {
	return s.data.CreatedOn
}

// ModifiedOn returns the last modification date of the data record.
func (s baseAttr) ModifiedOn() time.Time {
	return s.data.ModifiedOn
}

func (s *baseAttr) fillFromResponse(resp responseData) {
	s.data = resp
}
