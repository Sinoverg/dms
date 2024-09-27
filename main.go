package main

import (
	"log"
	"net/http"

	"github.com/Cirqach/dms/cmd/templ/body"
	"github.com/Cirqach/dms/cmd/templ/tables"
	"github.com/Cirqach/dms/internal/env"
	"github.com/a-h/templ"
)

func main() {
	if err := env.LoadEnvs(); err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}
	s := http.NewServeMux()
	s.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	s.Handle("/", templ.Handler(body.Body()))
	s.Handle("/templ/change-table", templ.Handler(tables.ChangeTablePage()))
	s.HandleFunc("/templ/add-button/{type}", addButtonHandler)
	http.ListenAndServe(":8080", s)
}

func addButtonHandler(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("type"){
	case 
	}
	err := tables.ChangeTablePage().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
