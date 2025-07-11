package mux

import (
	"net/http"
	"os"

	"github.com/danipurwadi/internal-transfer-system/app/middleware"
	"github.com/danipurwadi/internal-transfer-system/app/transferapp"
	"github.com/danipurwadi/internal-transfer-system/business/transferbus"
	"github.com/danipurwadi/internal-transfer-system/foundation/logger"
	"github.com/danipurwadi/internal-transfer-system/foundation/web"
)

func WebApi(log *logger.Logger, shutdown chan os.Signal) http.Handler {
	webClient := web.NewClient(shutdown, middleware.Logger(log), middleware.Errors(log))

	transferBus := transferbus.New()
	transferApp := transferapp.NewApp(transferBus)

	// register routes and handlers
	transferApp.Routes(webClient)
	return webClient
}
