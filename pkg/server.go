package server

import (
	"net/http"
)

type ServerInput struct {
	Port    string
	Routers *http.ServeMux
}

func NewServer(in ServerInput) error {
	err := http.ListenAndServe(in.Port, in.Routers)

	return err
}
