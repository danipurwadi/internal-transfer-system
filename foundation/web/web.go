// Package web contains a small web framework.
package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Handler represents a function that handles a http request
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Logger represents a function that will be called to add information
// to the logs.
type Logger func(ctx context.Context, msg string, v ...any)

type Client struct {
	*http.ServeMux
	mw []MidHandler
}

func NewClient(mw ...MidHandler) *Client {
	mux := http.NewServeMux()

	return &Client{
		ServeMux: mux,
		mw:       mw,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *Client) Handle(method string, path string, handler Handler, mw ...MidHandler) {
	// handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx := setValues(r.Context(), &v)

		if err := handler(ctx, w, r); err != nil {
			slog.Error("Failed to handle request", "err", err)
			return
		}
	}
	finalPath := fmt.Sprintf("%s %s", method, path)
	a.ServeMux.HandleFunc(finalPath, h)
}
