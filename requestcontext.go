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

func decorate(r *http.Request) *decoratedReadCloser {
	result := decoratedReadCloser{body: r.Body}
	r.Body = &result
	return &result
}

func (w *decoratedReadCloser) Read(p []byte) (n int, err error) {
	return w.body.Read(p)
}

func (w *decoratedReadCloser) Close() error {
	return w.body.Close()
}

func Set(r *http.Request, context context.Context) {
	drc := decoratedReadCloser{body: r.Body, context: context}
	r.Body = &drc
}

func Get(r *http.Request) context.Context {
	if c, ok := r.Body.(*decoratedReadCloser); ok {
		return c.context
	}
	return context.Background()
}
