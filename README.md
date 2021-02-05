# Form3 API Client Exercise

Author: Gabriel Nelle

## Considerations and Comments

### Core Data

Looking at the API documentation, I tried to figure out what would be the central information
for a developer using the client I was about to write. All the information that could be hidden
away, should be. Either hidden deeper or not exposed at all.

The attributes section seemed the central part. On creation and on fetching, there is more meta-data
that could be of interest to the user (developer). The tree above (data, list) does not need to be
exposed. Especially the list section is interesting for client developers, not so much for client users.

That's why I decided to make the attributes section central to my implementation. As input of the create
function and as output of `create`, `fetch` and `list`. As output the `Account` type not only contains the
account attributes but also provides functions to get to other meta-data like `ID`, `OrgID`, etc. That
data cannot be set though, e.g. when creating a new `Account` struct for creation.

### API Design

The API design is meant to be as simple as possible. To get a quick overview you can use `godoc` or `make doc`
and open `http://localhost:6060/pkg/github.com/tehsphinx/form3/` to see the documentation. Here a quick exammple:

```go
cl := form3.NewClient("http://localhost:8080")
account, err := cl.FetchAccount(ctx, accountID)
```

I was pondering a slightly different API-design for a while with a section per data type. It would have looked
this way from a users perspective:

```go
cl := form3.NewClient("http://localhost:8080")
account, err := cl.Account.Fetch(ctx, accountID)
```

For very large clients that might be a preferrable approach using autocomplete in the editor. I'd consider it
a matter of preference.

### Usability vs Performance

This client is mostly built around providing a simple, easy to use and hopefully intuitive API. Performance
was on my mind but not the first priority. For example the variadic option lists used in serveral places 
have the cost of extra heap allocations.

### Maintainability

Hiding away complexity to provide a simple API often comes at the cost of harder to maintain code. In this
client the complex function is `request`, which is able to handle all types of request encountered in this
test szenario. De-centralizing that into the individual handlers would have the disadvantage of many slightly
different implemmentations over time as many developers work on the code. The centralized approach keeps things
in one place and therefore one place to change. As complexity increases with more edge cases it might be good
to split it into multiple functions, e.g. one for `create`, `fetch`, `list`, `delete`, etc.

Another reason for the complex `request` function was for me to show some more complex code using some advanced
patterns. Especially in the `list` szenario, wrapping the meta-data into the `Account` struct was a bit tricky.

### Testing

Most tests are integratioon tests, testing the entire stack (against the server/database). They use the `form3_test` 
package to be able to use the API as a user would, restrained from using any internals.

Additionally, I added some unit tests for the validation function.

Running the tests is integrated in docker-compose:

```shell
docker-compose up
```

To enable colored communication logging: Set the `DEBUG` envorinment variable to `true` in the `docker-compose.yml`.

Also check `make list` for other commands.

### Dependencies

- **github.com/google/uuid**: UUID library for building uuids.
- **github.com/matryer/is**: Minimalistic assertion library used in the tests.
- **github.com/namsral/flag**: Drop in replacement for `flag` package that also reads environment variables. Used for reading environment variables
for running the tests.
- **github.com/tehsphinx/dbg**: Small debug library written by me (aged a bit / needs improvements). Used for colored output 
  in case debugging is enabled.
