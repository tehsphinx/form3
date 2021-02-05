package form3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

const accountsPath = "/v1/organisation/accounts"

// client side account validation errors
var (
	ErrInvalidCountry      = errors.New("country should match '^[A-Z]{2}$'")
	ErrInvalidBaseCurrency = errors.New("baseCurrency should match '^[A-Z]{3}$'")
	ErrInvalidBankID       = errors.New("bankID should match '^[A-Z0-9]{0,16}$'")
	ErrInvalidBankIDCode   = errors.New("bankIDCode should match '^[A-Z]{0,16}$'")
	ErrInvalidBIC          = errors.New("bic should match '^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$'")
	ErrInvalidAccClass     = errors.New("accountClassification should be one of [Personal Business]")
)

// Account holds account attributes.
type Account struct {
	baseAttr

	Country                 string   `json:"country"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BIC                     string   `json:"bic,omitempty"`
	IBAN                    string   `json:"iban,omitempty"`
	Name                    []string `json:"name,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	AccountClassification   string   `json:"account_classification,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  string   `json:"status,omitempty"`
	JointAccount            bool     `json:"joint_account,omitempty"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
	Switched                bool     `json:"switched,omitempty"`
}

// CreateAccount creates a new banking account.
func (s *Client) CreateAccount(ctx context.Context, orgID string, data *Account) (*Account, error) {
	if err := s.validateAccount(data); err != nil {
		return nil, fmt.Errorf("invalid Account information provided: %w", err)
	}

	uid := uuid.NewString()

	resp := &Account{}
	// resp := buildResp(respAttr)
	uri := s.buildURL(accountsPath, "", nil)
	if err := s.request(ctx, uri, typeAccounts, withMethod(http.MethodPost), withOrgID(orgID),
		withUID(uid), withReq(data), withResp(resp), withStatusOk(http.StatusCreated)); err != nil {
		return nil, err
	}

	return resp, nil
}

// FetchAccount retrieves the account information for given accound id.
func (s *Client) FetchAccount(ctx context.Context, uid string) (*Account, error) {
	resp := &Account{}
	uri := s.buildURL(accountsPath, uid, nil)
	if err := s.request(ctx, uri, typeAccounts, withResp(resp)); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Client) buildURL(basePath, uid string, params url.Values) string {
	if uid != "" {
		basePath += "/" + uid
	}

	uri := &url.URL{}
	uri.Path = basePath
	uri.RawQuery = params.Encode()
	return s.baseURL + uri.RequestURI()
}

// ListAccounts retrieves a list of accounts that can be filtered (not yet implemented) and has pagination.
func (s *Client) ListAccounts(ctx context.Context, opts ...ListOption) ([]Account, error) {
	params := url.Values{}
	applyOptions(params, opts)

	var accounts []Account
	uri := s.buildURL(accountsPath, "", params)
	if err := s.request(ctx, uri, typeAccounts,
		withListResp(
			func() responseFiller {
				return &Account{}
			},
			func(data responseFiller) {
				account := data.(*Account)
				accounts = append(accounts, *account)
			},
		),
	); err != nil {
		return nil, err
	}

	return accounts, nil
}

// DeleteAccount deletes the account with given account id. If the resource was not found a
// ErrNotFound will be returned. A ErrConflict indicates the account was updated meanwhile.
// Check for specific error types with e.g. `errors.Is(err, form3.ErrConflict)`
func (s *Client) DeleteAccount(ctx context.Context, uid string, version int) error {
	params := url.Values{}
	params.Set("version", strconv.Itoa(version))

	uri := s.buildURL(accountsPath, uid, params)
	return s.request(ctx, uri, typeAccounts,
		withMethod(http.MethodDelete), withStatusOk(http.StatusNoContent))
}

/* client side validation is not required as the server would deny an invalid request.
This catches problems without making a call. Best would be to import the validator
from the server code, so there is no extra effort to keep it in sync.
*/

// some client side validation. Does not need to be complete, but should never be stricter than server.
func getValidateAccount() func(attr *Account) error {
	countryRE := regexp.MustCompile("^[A-Z]{2}$")
	currencyRE := regexp.MustCompile("^[A-Z]{3}$")
	bankIDRE := regexp.MustCompile("^[A-Z0-9]{0,16}$")
	bankIDCodeRE := regexp.MustCompile("^[A-Z0-9]{0,16}$")
	bicRE := regexp.MustCompile("^([A-Z]{6}[A-Z0-9]{2}|[A-Z]{6}[A-Z0-9]{5})$")

	return func(attr *Account) error {
		switch {
		case !countryRE.MatchString(attr.Country):
			return ErrInvalidCountry
		case attr.BaseCurrency != "" && !currencyRE.MatchString(attr.BaseCurrency):
			return ErrInvalidBaseCurrency
		case attr.BankID != "" && !bankIDRE.MatchString(attr.BankID):
			return ErrInvalidBankID
		case attr.BankIDCode != "" && !bankIDCodeRE.MatchString(attr.BankIDCode):
			return ErrInvalidBankIDCode
		case attr.BIC != "" && !bicRE.MatchString(attr.BIC):
			return ErrInvalidBIC
		case attr.AccountClassification != "" && attr.AccountClassification != "Personal" &&
			attr.AccountClassification != "Business":
			return ErrInvalidAccClass
		}

		return nil
	}
}
