package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"tailscope/internal/app"
	"tailscope/internal/kube"
	"tailscope/internal/store"
)

type Server struct {
	app   *app.App
	store *store.Store
	kube  *kube.Client
}

func NewServer(
	application *app.App,
	kube *kube.Client,
) *Server {

	return &Server{
		app:   application,
		store: application.Store(),
		kube:  kube,
	}
}

func (s *Server) Router() http.Handler {

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{

		// React/Vite dev server
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://localhost:5175",
		},

		AllowedMethods: []string{
			"GET",
			"POST",
			"OPTIONS",
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},

		ExposedHeaders: []string{
			"Link",
		},

		AllowCredentials: true,

		MaxAge: 300,
	}))

	r.Get(
		"/api/logs",
		s.latestLogs,
	)

	r.Get(
		"/api/logs/search",
		s.searchLogs,
	)

	r.Get(
		"/api/deployments",
		s.deployments,
	)

	r.Get(
		"/api/logs/stream",
		s.streamLogs,
	)

	r.Get(
		"/api/namespaces",
		s.namespaces,
	)

	r.Get(
		"/api/pods",
		s.pods,
	)

	r.Post(
		"/api/logs/start",
		s.startLogs,
	)

	r.Post(
		"/api/logs/stop",
		s.stopLogs,
	)

	return r

}
