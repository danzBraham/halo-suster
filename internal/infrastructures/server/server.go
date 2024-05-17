package server

import (
	"log"
	"net/http"

	"github.com/danzBraham/halo-suster/internal/applications/services"
	"github.com/danzBraham/halo-suster/internal/helpers"
	user_repository_postgres "github.com/danzBraham/halo-suster/internal/infrastructures/repository/user"
	"github.com/danzBraham/halo-suster/internal/interfaces/http/api/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	Addr string
	DB   *pgxpool.Pool
}

func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{
		Addr: addr,
		DB:   db,
	}
}

func (s *APIServer) Launch() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		helpers.ResponseJSON(w, http.StatusOK, &helpers.ResponseBody{
			Message: "Welcome to Halo Suster API",
		})
	})

	userRepository := user_repository_postgres.NewUserRepositoryPostgres(s.DB)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/user", userController.Routes())
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		helpers.ResponseJSON(w, http.StatusNotFound, &helpers.ResponseBody{
			Message: "Route does not exist",
		})
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		helpers.ResponseJSON(w, http.StatusMethodNotAllowed, &helpers.ResponseBody{
			Message: "Method is not allowed",
		})
	})

	server := http.Server{
		Addr:    s.Addr,
		Handler: r,
	}

	log.Printf("Server listening on %s\n", s.Addr)
	return server.ListenAndServe()
}
