package users

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/latdev/httpxd/system/syscore"
)

type Handler struct {
	*syscore.CoreDirector
}

func (h *Handler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	router := mux.NewRouter()
	router.NotFoundHandler = h.CoreDirector.MuxServeNotFound()
	router.HandleFunc("/auth", h.handleAuthenticateUser).Methods("POST")
	router.ServeHTTP(wr, req)
}

func (h *Handler) handleAuthenticateUser(wr http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(wr, "OK")
}

