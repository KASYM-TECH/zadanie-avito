package server

import (
	"avito/controllers"
	"avito/log"
	"github.com/julienschmidt/httprouter"
)

type Router struct {
	Router *httprouter.Router
	logger log.Logger
}

func NewRouter(logger log.Logger) *Router {
	return &Router{Router: httprouter.New(), logger: logger}
}

func (r *Router) AddRoutes(m *Middleware,
	userController *controllers.UserController,
	bannerController *controllers.BannerController) {

	r.Router.POST("/auth/signup", m.Wrap(userController.Signup))
	r.Router.POST("/auth/login", m.Wrap(userController.Login))

	r.Router.GET("/banner", m.WrapAuthed(bannerController.GetBanner))
}
