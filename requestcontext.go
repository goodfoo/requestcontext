// Package requestcontext is a contention free http.Request adapter for golang.org/x/net/context
package requestcontext

import (
	"io"
	"net/http"

	"golang.org/x/net/context"
)

type decoratedReadCloser struct {
	body    io.ReadCloser
	context context.Context
}

func (w *decoratedReadCloser) Read(p []byte) (n int, err error) {
	return w.body.Read(p)
}

func (w *decoratedReadCloser) Close() error {
	return w.body.Close()
}

// Set this context
// call Get() prior and wrap the currency context with
// one of the golang.org/x/net/context methods - see example
func Set(r *http.Request, context context.Context) {
	drc := decoratedReadCloser{body: r.Body, context: context}
	r.Body = &drc
}

// Get returns a parent context, will make one if needed.
// always call this before Set() to preserve context chain
func Get(r *http.Request) context.Context {
	if c, ok := r.Body.(*decoratedReadCloser); ok {
		return c.context
	}
	return context.Background()
}
