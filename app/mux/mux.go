package mux

import (
	"net/http"
	"os"

	"github.com/danipurwadi/internal-transfer-system/app/middleware"
	"github.com/danipurwadi/internal-transfer-system/app/transactionapp"
	"github.com/danipurwadi/internal-transfer-system/foundation/logger"
	"github.com/danipurwadi/internal-transfer-system/foundation/web"
)

func WebApi(log *logger.Logger, shutdown chan os.Signal) http.Handler {
	app := web.NewApp(shutdown, middleware.Logger(log), middleware.Errors(log))

	transactionapp.Routes(app)
	return app
}
