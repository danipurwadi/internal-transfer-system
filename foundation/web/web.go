package web

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// Handler represents a function that handles a http request
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Logger represents a function that will be called to add information
// to the logs.
type Logger func(ctx context.Context, msg string, v ...any)

type App struct {
	// log Logger
	*http.ServeMux
	shutdown chan os.Signal
	mw       []MidHandler
}

func NewApp(shutdown chan os.Signal, mw ...MidHandler) *App {
	mux := http.NewServeMux()

	return &App{
		ServeMux: mux,
		shutdown: shutdown,
		mw:       mw,
		// log:      log,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) HandleFunc(pattern string, handler Handler, mw ...MidHandler) {
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

	a.ServeMux.HandleFunc(pattern, h)
}
