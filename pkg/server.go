package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ServerInput struct {
	Port    string
	Routers *mux.Router
}

func NewServer(in ServerInput) error {
	err := http.ListenAndServe(in.Port, in.Routers)

	return err
}
