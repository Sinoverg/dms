package main

import (
	"log"
	"net/http"

	"github.com/Cirqach/dms/cmd/handler"
	"github.com/Cirqach/dms/cmd/templ/body"
	"github.com/Cirqach/dms/cmd/templ/tables"
	"github.com/Cirqach/dms/internal/database"
	"github.com/Cirqach/dms/internal/env"
	"github.com/a-h/templ"
)

func main() {
	if err := env.LoadEnvs(); err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}
	db := database.NewController()
	s := http.NewServeMux()
	s.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	s.Handle("/", templ.Handler(body.Body()))
	s.Handle("/templ/change-table", templ.Handler(tables.ChangeTablePage()))
	s.HandleFunc("/templ/tables/{table}", handler.TableHandler(db))
	s.HandleFunc("/templ/buttons/{Button}/{Type}", handler.ButtonHandler)
	s.HandleFunc("/api/{function}/{table}", handler.ApiHandler(db))
	log.Println("Running server on port 8080")
	log.Fatal(http.ListenAndServe(":1337", s))
	// http.ListenAndServeTLS(":80", "ss-cert.crt", "private-key.key", s)
}
