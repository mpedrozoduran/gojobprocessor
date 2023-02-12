package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	MuxRouter *chi.Mux
}

func (r *Router) InitApiJobs() {
	r.MuxRouter.Use(middleware.Logger)

	r.MuxRouter.Route("/jobprocessor", func(rtr chi.Router) {
		rtr.Get("/jobs", GetJob)
		rtr.Get("/jobs/all", GetJobs)
		rtr.Post("/jobs/create", CreateJob)
		rtr.Post("/jobs/run", RunJob)
	})
}
