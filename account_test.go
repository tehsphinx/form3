package form3

import (
	"testing"

	"github.com/matryer/is"
)

func Test_getValidateAccount(t *testing.T) {
	tests := []struct {
		name    string
		attr    *Account
		wantErr error
	}{
		{
			name: "invalid country",
			attr: &Account{
				Country: "G",
			},
			wantErr: ErrInvalidCountry,
		},
		{
			name: "invalid currency",
			attr: &Account{
				Country:      "GB",
				BaseCurrency: "GB",
			},
			wantErr: ErrInvalidBaseCurrency,
		},
		{
			name: "invalid bank id",
			attr: &Account{
				Country: "GB",
				BankID:  "2fdsjo",
			},
			wantErr: ErrInvalidBankID,
		},
		{
			name: "invalid bank id code",
			attr: &Account{
				Country:    "GB",
				BankIDCode: "fjio23",
			},
			wantErr: ErrInvalidBankIDCode,
		},
		{
			name: "invalid BIC",
			attr: &Account{
				Country: "GB",
				BIC:     "43278r23",
			},
			wantErr: ErrInvalidBIC,
		},
		{
			name: "invalid account classification",
			attr: &Account{
				Country:               "GB",
				AccountClassification: "fjwi",
			},
			wantErr: ErrInvalidAccClass,
		},
		{
			name: "valid account",
			attr: &Account{
				Country:               "GB",
				BaseCurrency:          "GBP",
				BankID:                "400300",
				BankIDCode:            "GBDSC",
				BIC:                   "NWBKGB22",
				AccountClassification: "Personal",
			},
			wantErr: nil,
		},
	}

	validateAccount := getValidateAccount()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := is.New(t)

			err := validateAccount(tt.attr)
			assert.Equal(err, tt.wantErr)
		})
	}
}
