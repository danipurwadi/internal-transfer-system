package transferapp

import (
	"context"
	"net/http"

	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/foundation/web"
)

type app struct {
	transferbus *transferbus.Bus
}

func NewApp(bus *transferbus.Bus) *app {
	return &app{
		transferbus: bus,
	}
}

func (a *app) Routes(mux *web.Client) {
	mux.HandleFunc("GET /health", a.health)
}

func (a *app) health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status bool `json:"status"`
	}{
		Status: true,
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
