package route

import (
	"microservice/internal/controller"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Router struct {
	Logger        *zap.Logger
	Router        *chi.Mux
	ControllerID  controller.HTTPHandler
	ControllerWeb controller.HTTPHandler
}

func NewRouter(logger *zap.Logger, controllerID controller.HTTPHandler, controllerWeb controller.HTTPHandler) *Router {

	return &Router{
		Logger:        logger,
		Router:        chi.NewRouter(),
		ControllerID:  controllerID,
		ControllerWeb: controllerWeb,
	}
}

func (r *Router) MyRouter() {
	r.Logger.Info("add road /id/{id}")
	r.Router.Get("/id/{id}", r.ControllerID.ServeHTTP)

	r.Router.Get("/index.html", r.ControllerWeb.ServeHTTP)
}
