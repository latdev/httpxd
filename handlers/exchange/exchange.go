package exchange

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/latdev/httpxd/system/syscore"
	"net/http"
)

type Handler struct {
	*syscore.CoreDirector
}

func (h *Handler) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	router := mux.NewRouter()
	router.NotFoundHandler = h.CoreDirector.MuxServeNotFound()
	router.HandleFunc("/{from}/{targ}/{amount:[0-9\\.]+}", h.handleRateExchange).Methods("GET")
	router.ServeHTTP(wr, req)
}

func (h *Handler) handleRateExchange(wr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	session := h.Session(wr, req)
	defer session.Save()

	var n int = 0
	if q, ok := session.Get("number"); ok {
		if qn, ok := q.(int); ok {
			n = qn
		}
	}
	n += 1
	session.Set("number", n)

	fmt.Fprintf(wr, "from %s to %s %s and N=%d", vars["from"], vars["targ"], vars["amount"], n)
}