package server

import (
	"log"
	"net/http"

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
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Halo Suster API\n"))
	})

	server := http.Server{
		Addr:    s.Addr,
		Handler: mux,
	}

	log.Printf("Server listening on %s\n", s.Addr)
	return server.ListenAndServe()
}
