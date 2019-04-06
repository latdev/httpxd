package main

import (
	"fmt"
	"github.com/latdev/httpxd/handlers/exchange"
	"github.com/latdev/httpxd/system/syslogger"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/latdev/httpxd/system/syscore"
	"github.com/latdev/httpxd/handlers/users"
)

func main() {

	if director, err := syscore.New(); err == nil {
		defer director.Close()

		router := mux.NewRouter().StrictSlash(true)
		router.NotFoundHandler = director.MuxServeNotFound()
		router.HandleFunc("/favicon.ico", func(wr http.ResponseWriter, _ *http.Request) {
			wr.Header().Set("Content-Type", "image/x-icon")
			wr.Header().Set("Cache-Control", "public, max-age=7776000")
			fmt.Fprintln(wr, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=\n")
		}).Methods("GET", "HEAD")
		router.PathPrefix("/io/exchange/").Handler(http.StripPrefix("/io/exchange", &exchange.Handler{CoreDirector: director}))
		router.PathPrefix("/io/user/").Handler(http.StripPrefix("/io/user", &users.Handler{CoreDirector: director}))

		server := &http.Server{
			Addr:           director.Settings.Server.Binding,
			Handler:        &syslogger.Handler{Next: router},
			ReadTimeout:    time.Duration(director.Settings.Server.ReadTimeout) * time.Second,
			WriteTimeout:   time.Duration(director.Settings.Server.WriteTimeout) * time.Second,
			MaxHeaderBytes: int(director.Settings.Server.MaxHeaderBytes),
		}

		log.Printf("starting server on %s", server.Addr)
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(err)
	}
}
