package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/sidgim/finance_user/internal/user"
	"net/http"
)

func NewRouter(
	uH *user.Handler,
) http.Handler {
	r := chi.NewRouter()
	r.Route("/users", uH.Mount)
	return r
}
