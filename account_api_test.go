package form3_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/namsral/flag"
	"github.com/tehsphinx/form3"
)

var endpoint string
var debugEnabled bool

const orgID = "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"

func TestMain(m *testing.M) {
	flag.StringVar(&endpoint, "endpoint", "http://localhost:8080", "test server endpoint url")
	flag.BoolVar(&debugEnabled, "debug", false, "enable colored debug output")
	flag.Parse()

	cleanAccountsTable()

	code := m.Run()
	os.Exit(code)
}

func getClient() *form3.Client {
	var options []form3.ClientOption
	if debugEnabled {
		options = append(options, form3.WithDebug())
	}
	cl := form3.NewClient(endpoint, options...)
	return cl
}

/*
Since we want to test the client, not the server, there is little point testing all the error szenarios of the server.
The validation will be tested separately without making calls to the server.

Go tests are executed sequentially and in order by default. I'm using that here to keep the tests simple.
If I wanted to execute them in parallell, the tests would need to be written differently,
each writing their own data, before fetching and deleting it again. Since this is running against a database,
the test data would have to be without overlap. E.g the `List` call would need a filter that would exclude all
other test data or its own server/database to run against.

WARNING: the tests will clean the Account table first.
*/

var accountTests = []struct {
	orgID       string
	uid         string
	version     int
	name        string
	createData  *form3.Account
	accountData *form3.Account
}{
	{
		name:  "UK account without CoP",
		orgID: orgID,
		createData: &form3.Account{
			Country:      "GB",
			BaseCurrency: "GBP",
			BankID:       "400300",
			BankIDCode:   "GBDSC",
			BIC:          "NWBKGB22",
		},
		accountData: &form3.Account{
			Country:      "GB",
			BaseCurrency: "GBP",
			BankID:       "400300",
			BankIDCode:   "GBDSC",
			BIC:          "NWBKGB22",
		},
	},
	{
		name:  "UK account with CoP",
		orgID: orgID,
		createData: &form3.Account{
			Country:                 "GB",
			BaseCurrency:            "GBP",
			BankID:                  "400300",
			BankIDCode:              "GBDSC",
			BIC:                     "NWBKGB22",
			Name:                    []string{"Samantha Holder"},
			AlternativeNames:        []string{"Sam Holder"},
			AccountClassification:   "Personal",
			JointAccount:            false,
			AccountMatchingOptOut:   false,
			SecondaryIdentification: "A1B2C3D4",
		},
		accountData: &form3.Account{
			Country:                 "GB",
			BaseCurrency:            "GBP",
			BankID:                  "400300",
			BankIDCode:              "GBDSC",
			BIC:                     "NWBKGB22",
			AccountClassification:   "Personal",
			JointAccount:            false,
			AccountMatchingOptOut:   false,
			SecondaryIdentification: "A1B2C3D4",
		},
	},
}

func TestClient_CreateAccount(t *testing.T) {
	for i, tt := range accountTests {
		i := i // scope variable
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			cl := getClient()
			got, err := cl.CreateAccount(ctx, tt.orgID, tt.createData)

			assert := is.New(t)
			assert.NoErr(err)
			assert.True(got.ID() != "")
			assert.Equal(copyAccount(*got), tt.accountData)

			// save uid to test case for the other tests
			accountTests[i].uid = got.ID()
			accountTests[i].version = got.Version()
		})
	}
}

func TestClient_FetchAccount(t *testing.T) {
	for _, tt := range accountTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			cl := getClient()
			got, err := cl.FetchAccount(ctx, tt.uid)

			assert := is.New(t)
			assert.NoErr(err)
			assert.True(got.ID() != "")
			assert.Equal(copyAccount(*got), tt.accountData)
		})
	}
}

func TestClient_ListAccounts(t *testing.T) {
	tests := []struct {
		name         string
		options      []form3.ListOption
		expectedUIDs []string
	}{
		{
			name:         "no pagination",
			expectedUIDs: []string{accountTests[0].uid, accountTests[1].uid},
		},
		{
			name: "with page option",
			options: []form3.ListOption{
				form3.WithPageNo(0),
			},
			expectedUIDs: []string{accountTests[0].uid, accountTests[1].uid},
		},
		{
			name: "with size option",
			options: []form3.ListOption{
				form3.WithPageSize(5),
			},
			expectedUIDs: []string{accountTests[0].uid, accountTests[1].uid},
		},
		{
			name: "second page has no data",
			options: []form3.ListOption{
				form3.WithPageNo(1),
				form3.WithPageSize(5),
			},
			expectedUIDs: []string{},
		},
		{
			name: "first page with size 1",
			options: []form3.ListOption{
				form3.WithPageNo(0),
				form3.WithPageSize(1),
			},
			expectedUIDs: []string{accountTests[0].uid},
		},
		{
			name: "second page with size 1",
			options: []form3.ListOption{
				form3.WithPageNo(1),
				form3.WithPageSize(1),
			},
			expectedUIDs: []string{accountTests[1].uid},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := is.New(t)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			cl := getClient()
			got, err := cl.ListAccounts(ctx, tt.options...)

			// assert that there was no error and the expected amount of results
			assert.NoErr(err)
			assert.Equal(len(got), len(tt.expectedUIDs))

			// check if the expected uids are in the result
			for _, uid := range tt.expectedUIDs {
				var found bool
				for _, account := range got {
					assert.True(account.ID() != "")
					if uid != account.ID() {
						continue
					}
					found = true

				}
				assert.True(found)
			}

			// check if the results are the same as expected in accountTests
			for _, test := range accountTests {
				for _, account := range got {
					if test.uid != account.ID() {
						continue
					}
					assert.Equal(copyAccount(account), test.accountData)
				}
			}
		})
	}
}

func TestClient_DeleteAccount(t *testing.T) {
	for _, tt := range accountTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			assert := is.NewRelaxed(t)
			cl := getClient()

			// try to delete wrong version No
			err := cl.DeleteAccount(ctx, tt.uid, tt.version+1)
			// Note: seems the way the API is implemented server side, a wrong version returns a 404 if that
			// version does not exist at all, not a 409 as one could have read from the documentation.
			// Is not really a use case that matters though, since normally one would not invent a version number.
			assert.True(errors.Is(err, form3.ErrNotFound))

			// test deletion
			err = cl.DeleteAccount(ctx, tt.uid, tt.version)
			assert.NoErr(err)
		})
	}
}

// use to get rid of private values for comparison
func copyAccount(account form3.Account) *form3.Account {
	return &form3.Account{
		Country:                 account.Country,
		BaseCurrency:            account.BaseCurrency,
		AccountNumber:           account.AccountNumber,
		BankID:                  account.BankID,
		BankIDCode:              account.BankIDCode,
		BIC:                     account.BIC,
		IBAN:                    account.IBAN,
		Name:                    account.Name,
		AlternativeNames:        account.AlternativeNames,
		AccountClassification:   account.AccountClassification,
		JointAccount:            account.JointAccount,
		AccountMatchingOptOut:   account.AccountMatchingOptOut,
		SecondaryIdentification: account.SecondaryIdentification,
		Switched:                account.Switched,
		Status:                  account.Status,
	}
}

func cleanAccountsTable() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cl := form3.NewClient(endpoint)

	accounts, err := cl.ListAccounts(ctx, form3.WithPageSize(100))
	if err != nil {
		log.Fatal(err)
	}
	if len(accounts) == 100 {
		log.Fatal("tests would delete too much data. aborting")
	}

	for _, account := range accounts {
		if err := cl.DeleteAccount(ctx, account.ID(), account.Version()); err != nil {
			log.Fatal(err)
		}
	}
}
