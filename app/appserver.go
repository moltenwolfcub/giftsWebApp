package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/moltenwolfcub/giftsWebApp/config"
	"github.com/moltenwolfcub/giftsWebApp/controller"
)

type AppServer struct {
	cfg    config.Config
	router http.Handler
}

func NewAppServer(db *sql.DB, cfg *config.Config) *AppServer {
	return &AppServer{
		router: controller.BuildRouter(db),
		cfg:    *cfg,
	}
}

func (a *AppServer) Run() {
	log.Printf("Starting server on :%d\n", a.cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.Port), a.router))
}
