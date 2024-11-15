package main

import (
	"net/http"

	"github.com/Cirqach/dms/cmd/handler"
	"github.com/Cirqach/dms/cmd/templ/body"
	"github.com/Cirqach/dms/cmd/templ/tables"
	"github.com/Cirqach/dms/internal/database"
	"github.com/Cirqach/dms/internal/env"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := env.LoadEnvs(); err != nil {
		log.Fatal().Err(err).Msg("Error due loading environment variables")
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	db,err := database.NewController()
	if err != nil {
		log.Fatal().Err(err).Msg("Error due creating new controller")
	}
	HandleAPI(db,r)
	HandleTempl(db,r)
	r.Handle("/static/*", http.StripPrefix("/static/",http.FileServer(http.Dir("./static/"))))
	log.Info().Msg("Running server on port 1337")
	log.Fatal().Err(http.ListenAndServe(":1337", r)).Msg("Error listen and serve")
	// http.ListenAndServeTLS(":1337", "ss-cert.crt", "private-key.key", r)
}

func HandleTempl(db *database.DBController,r *chi.Mux){
	r.Handle("/", templ.Handler(body.Body()))
	r.Handle("/templ/change-table", templ.Handler(tables.ChangeTablePage()))
	r.HandleFunc("/templ/tables/{table}", handler.TableHandler(db))
}

func HandleAPI(db *database.DBController, r *chi.Mux) {
	r.HandleFunc("/api/delete/{table}",handler.ApiDelete(db))

	r.HandleFunc("/api/add/users",handler.ApiAddUser(db))
	r.HandleFunc("/api/add/videofiles",handler.ApiAddVideofile(db))
	r.HandleFunc("/api/add/broadcasts",handler.ApiAddBroadcast(db))
	r.HandleFunc("/api/add/broadcasts_files",handler.ApiAddBroadcastsFiles(db))
	r.HandleFunc("/api/add/broadcasts_users",handler.ApiAddBroadcastsUsers(db))

	r.HandleFunc("/api/update/videofiles",handler.ApiUpdateVideofile(db))
	r.HandleFunc("/api/update/users",handler.ApiUpdateUser(db))
	r.HandleFunc("/api/update/broadcasts",handler.ApiUpdateBroadcast(db))

} 