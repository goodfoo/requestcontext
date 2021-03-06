package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/goodfoo/requestcontext"
	"golang.org/x/net/context"
)

// arbitrary value type
type blah struct {
	s    string
	rune rune
	out  int
}

// any comparable
var comparable = "key"

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
	c := context.WithValue(parent, comparable, value)

	// no contention here ever
	requestcontext.Set(r, c)

	// good luck little buddy
	next(rw, r)
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
	// get the context - no contention here
	c := requestcontext.Get(r)

	// could call Deadline(), or Done() or Err()
	// call Value and convert the response to arbitrary type
	if b, ok := c.Value(comparable).(blah); ok {
		fmt.Fprintf(w, "Context value %s %c %d!", b.s, b.rune, b.out)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", MyHandler)

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(MyMiddleware))
	n.UseHandler(mux)
	n.Run(":3000")
}
