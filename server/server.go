package server

import (
	"net/http"

	"github.com/microphoneabuser/comics_service/pkg/handler"
)

func RunServer(port string, handlers *handler.Handler) error {
	handlers.SetupRoutes()
	return http.ListenAndServe(":"+port, nil)
}
