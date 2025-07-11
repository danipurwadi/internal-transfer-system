package transactionapp

import (
	"context"
	"net/http"

	"github.com/danipurwadi/internal-transfer-system/foundation/web"
)

func Routes(mux *web.App) {
	mux.HandleFunc("GET /health", health)
}

func health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status bool `json:"status"`
	}{
		Status: true,
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
