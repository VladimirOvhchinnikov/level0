package route

import (
	"microservice/internal/controller"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Router struct {
	Logger     *zap.Logger
	Router     *chi.Mux
	Controller controller.HTTPHandler
}

func NewRouter(logger *zap.Logger, controller controller.HTTPHandler) *Router {

	return &Router{
		Logger:     logger,
		Router:     chi.NewRouter(),
		Controller: controller,
	}
}

func (r *Router) MyRouter() {
	r.Logger.Info("add road /id/{id}")
	r.Router.Get("/id/{id}", r.Controller.ServeHTTP)
}
