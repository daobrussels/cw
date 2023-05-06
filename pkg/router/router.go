package router

import (
	"fmt"
	"net/http"

	"github.com/daobrussels/cw/pkg/config"
	"github.com/daobrussels/cw/pkg/push"
	"github.com/daobrussels/cw/pkg/server"
	"github.com/daobrussels/cw/pkg/token"
	"github.com/daobrussels/cw/pkg/transaction"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	// ...
	conf *config.Config
}

func NewServer(conf *config.Config) server.Server {
	return &Router{
		conf: conf,
	}
}

// implement the Server interface
func (r *Router) Start(port int) error {

	cr := chi.NewRouter()

	// configure middleware
	cr.Use(middleware.Compress(5))
	cr.Use(HealthMiddleware)
	cr.Use(SignatureMiddleware)

	// instantiate handlers
	transaction := transaction.NewHandlers()
	token := token.NewHandlers()
	push := push.NewHandlers()

	// configure routes
	cr.Post("/transaction", transaction.Send)

	cr.Route("/token", func(cr chi.Router) {
		cr.Post("mint", token.Mint)
		cr.Post("/burn", token.Burn)
	})

	cr.Route("/push", func(cr chi.Router) {
		cr.Put("/associate", push.Associate)
		cr.Delete("/dissociate", push.Dissociate)
	})

	// start the server
	return http.ListenAndServe(fmt.Sprintf(":%v", port), cr)
}
