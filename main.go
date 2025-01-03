package main

import (
	"net/http"

	"github.com/Cirqach/dms/cmd/handler"
	"github.com/Cirqach/dms/cmd/templ/body"
	"github.com/Cirqach/dms/cmd/templ/tables"
	"github.com/Cirqach/dms/internal/database"
	"github.com/Cirqach/dms/internal/database/models"
	"github.com/Cirqach/dms/internal/env"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TODO:
// 1. Change video file user id input to select.
// 2. Search
// 3. Filtration
// 4. Sorting
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if err := env.LoadEnvs(); err != nil {
		log.Fatal().Err(err).Msg("Error due loading environment variables")
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	db, err := database.NewController()
	if err != nil {
		log.Fatal().Err(err).Msg("Error due creating new controller")
	}
	HandleAPI(db, r)
	HandleTempl(db, r)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Info().Msg("Running server on port 1488")
	log.Fatal().Err(http.ListenAndServe(":1488", r)).Msg("Error listen and serve")
	// http.ListenAndServeTLS(":1337", "ss-cert.crt", "private-key.key", r)
}

func HandleTempl(db *database.DBController, r *chi.Mux) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.SelectAll("users")
		if err != nil {
			log.Err(err).Msg("Error selecting all records from users table")
			http.Error(w, "Internal error, try later", http.StatusInternalServerError)
			return
		}
		users := make([]models.User, 0)
		for rows.Next() {
			user := models.User{}
			err := rows.Scan(&user.Id, &user.Fname, &user.Sname, &user.Nickname, &user.Login, &user.Email, &user.Password)
			if err != nil {
				log.Err(err).Msg("Error scanning users rows")
				http.Error(w, "Internal error, try later", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}
		rows, err = db.SelectAll("broadcasts")
		if err != nil {
			log.Err(err).Msg("Error selecting all records from videofiles table")
			http.Error(w, "Internal error, try later", http.StatusInternalServerError)
			return
		}
		broadcasts := make([]models.Broadcast, 0)
		for rows.Next() {
			b := models.Broadcast{}
			err := rows.Scan(&b.Id, &b.StartTime, &b.EndTime)
			if err != nil {
				log.Err(err).Msg("Error scanning broadcasts rows")
				http.Error(w, "Internal error, try later", http.StatusInternalServerError)
				return
			}
			broadcasts = append(broadcasts, b)
		}
		rows, err = db.SelectAll("videofiles")
		if err != nil {
			log.Err(err).Msg("Error selecting all records from videofiles table")
			http.Error(w, "Internal error, try later", http.StatusInternalServerError)
			return
		}
		videofiles := make([]models.Videofile, 0)
		for rows.Next() {
			v := models.Videofile{}
			err := rows.Scan(&v.Id, &v.Filename, &v.Uploader, &v.Size, &v.Duration)
			if err != nil {
				log.Err(err).Msg("Error scanning videofiles rows")
				http.Error(w, "Internal error, try later", http.StatusInternalServerError)
				return
			}
			videofiles = append(videofiles, v)
		}

		body.Body(users, broadcasts, videofiles).Render(r.Context(), w)
	})
	r.Handle("/templ/change-table", templ.Handler(tables.ChangeTablePage()))
	r.HandleFunc("/templ/tables/{table}", handler.TableHandler(db))
	r.HandleFunc("/search/{table}", handler.Search(db))
	r.HandleFunc("/filter/{table}", handler.Filter(db))
	r.HandleFunc("/sort/{table}/{type}/{order}", handler.Sort(db))

}

func HandleAPI(db *database.DBController, r *chi.Mux) {
	r.HandleFunc("/api/delete/{table}/{id}", handler.ApiDelete(db))

	r.HandleFunc("/api/add/users", handler.ApiAddUser(db))
	r.HandleFunc("/api/add/videofiles", handler.ApiAddVideofile(db))
	r.HandleFunc("/api/add/broadcasts", handler.ApiAddBroadcast(db))
	r.HandleFunc("/api/add/broadcasts-files", handler.ApiAddBroadcastsFiles(db))
	r.HandleFunc("/api/add/broadcasts-users", handler.ApiAddBroadcastsUsers(db))

	r.HandleFunc("/api/patch/videofiles/{id}", handler.ApiUpdateVideofile(db))
	r.HandleFunc("/api/patch/users/{id}", handler.ApiUpdateUser(db))
	r.HandleFunc("/api/patch/broadcasts/{id}", handler.ApiUpdateBroadcast(db))
	r.HandleFunc("/api/patch/broadcasts-files/{id}", handler.ApiUpdateBroadcastsFiles(db))
	r.HandleFunc("/api/patch/broadcasts-users/{id}", handler.ApiUpdateBroadcastUsers(db))

}
