# requestcontext

Package requestcontext is a contention free http.Request adapter for golang.org/x/net/context for [Go](https://golang.org/).

Handy for passing request level context through middleware to your handlers, plus all the other context goodness.

Our team was frustrated by contention related performance issues with other context packages.

Motivated by
[Go context package](https://godoc.org/golang.org/x/net/context) and
[James Abley Gorilla Context pull request](https://github.com/gorilla/context/pull/21)

Not limited to string -> string !

Use any comparable as a key and any value type.

Leverage all the awesomeness of the golang.org/x/net/context using traditional http.Handler's

## Usage

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/goodfoo/requestcontext"
	"golang.org/x/net/context"
)

type blah struct {
	s    string
	rune rune
	out  int
}

// MyMiddleware assigns a value into the request
// any number of middleware may be nested
// it could register a Cancel func(), a Timeout or Deadline
func MyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// arbitrary value type
	value := blah{
		"a☻☺b", '☺', 4,
	}

	// always Get the context, will be populated if needed
	parent := requestcontext.Get(r)

	// this could be any of the google context extensions
	// WithCancel, WithDeadline, WithTimeout,
	// or the ever popular WithValue
	c := context.WithValue(parent, "key", value)

	// no contention here ever
	requestcontext.Set(r, c)

	// good luck little buddy
	next(rw, r)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get the context - no contention here
		c := requestcontext.Get(r)

		// could call Deadline(), or Done() or Err()
		// convert the response to arbitrary type
		if b, ok := c.Value("key").(blah); ok {
			fmt.Fprintf(w, "Context value %s %c %d!", b.s, b.rune, b.out)
		}
	})

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(MyMiddleware))
	n.UseHandler(mux)
	n.Run(":3000")
}

```