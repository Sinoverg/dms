package main

import (
	"log"
	"net/http"

	"github.com/Cirqach/dms/cmd/handler"
	"github.com/Cirqach/dms/cmd/templ/body"
	"github.com/Cirqach/dms/cmd/templ/buttons"
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
	s.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	s.Handle("/", templ.Handler(body.Body()))
	s.Handle("/templ/change-table", templ.Handler(tables.ChangeTablePage()))
	s.Handle("/tables/users", templ.Handler(handler.TableHandler(db, "users")))
	s.Handle("/tables/files", templ.Handler(handler.TableHandler(db, "files")))
	s.Handle("/tables/broadcast", templ.Handler(handler.TableHandler(db, "broadcast")))
	s.HandleFunc("/templ/buttons/{button}/{type}", addButtonHandler)
	http.ListenAndServe(":8080", s)
}

func addButtonHandler(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("button") {
	case "ask":
		switch r.PathValue("type") {
		case "ask":
			buttons.AddAskTable().Render(r.Context(), w)
		case "user":
			buttons.AddUser().Render(r.Context(), w)
		case "file":
			buttons.AddFile().Render(r.Context(), w)
		case "broadcast":
			buttons.AddBroadcast().Render(r.Context(), w)
		}

	}
	buttons.AddAskTable().Render(r.Context(), w)
}
