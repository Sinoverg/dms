package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Cirqach/dms/cmd/templ/messages"
	"github.com/Cirqach/dms/cmd/templ/tables"
	"github.com/Cirqach/dms/internal/database"
	"github.com/Cirqach/dms/internal/database/models"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {

}

func TableHandler(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get request: ", r.URL.Path)
		table := r.PathValue("table")
		rows, err := db.SelectAll(table)
		if err != nil {
			messages.Error(err).Render(r.Context(), w)
			return
		}
		switch table {
		case "videofiles":
			files := make([]models.Videofile, 0)
			for rows.Next() {
				var id, filename, uploader, size, duration string
				err := rows.Scan(&id, &filename, &uploader, &size, &duration)
				if err != nil {
					messages.Error(err).Render(r.Context(), w)
					return
				}
				file := models.Videofile{
					Id:       id,
					Filename: filename,
					Uploader: uploader,
					Size:     size,
					Duration: duration,
				}
				files = append(files, file)
			}
			tables.Files(files).Render(r.Context(), w)
			return
		case "users":
			users := make([]models.User, 0)
			for rows.Next() {
				var id, fname, sname, nickname, login, email, password string
				err := rows.Scan(&id, &fname, &sname, &nickname, &login, &email, &password)
				if err != nil {
					messages.Error(err).Render(r.Context(), w)
					return
				}
				user := models.User{
					Id:       id,
					Fname:    fname,
					Sname:    sname,
					Nickname: nickname,
					Login:    login,
					Email:    email,
					Password: password,
				}
				users = append(users, user)
			}
			tables.Users(users).Render(r.Context(), w)
			return
		case "broadcasts":
			broadcasts := make([]models.Broadcast, 0)
			for rows.Next() {
				var id string
				var starttime, endtime time.Time
				err := rows.Scan(&id, &starttime, &endtime)
				if err != nil {
					messages.Error(err).Render(r.Context(), w)
					return
				}
				broadcast := models.Broadcast{
					Id:                 id,
					BroadcastStartTime: starttime,
					BroadcastEndTime:   endtime,
				}
				broadcasts = append(broadcasts, broadcast)
			}
			tables.Broadcast(broadcasts).Render(r.Context(), w)
			return
		case "broadcasts_users":
			broadcastsUsers := make([]models.BroadcastUsers, 0)
			for rows.Next() {
				var broadcastid, userid string
				err := rows.Scan(&broadcastid, &userid)
				if err != nil {
					messages.Error(err).Render(r.Context(), w)
					return
				}
				broadcast := models.BroadcastUsers{
					BroadcastId: broadcastid,
					UserId:      userid,
				}
				broadcastsUsers = append(broadcastsUsers, broadcast)
			}
			tables.Broadcasts_users(broadcastsUsers).Render(r.Context(), w)
			return
		case "broadcasts_files":
			broadcastFiles := make([]models.BroadcastFiles, 0)
			for rows.Next() {
				var broadcastid, fileid string
				err := rows.Scan(&broadcastid, &fileid)
				if err != nil {
					messages.Error(err).Render(r.Context(), w)
					return
				}
				bf := models.BroadcastFiles{
					BroadcastId: broadcastid,
					VideofileId: fileid,
				}
				broadcastFiles = append(broadcastFiles, bf)
			}
			tables.Broadcasts_files(broadcastFiles).Render(r.Context(), w)
			return
		}
	}
}

func ApiHandler(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Take %s request from %s\n", r.Method, r.URL)
		function := r.PathValue("function")
		table := r.PathValue("table")
		switch function {
		case "add":
			var msg string
			switch table {
			case "users":
				fname := r.FormValue("first_name")
				sname := r.FormValue("second_name")
				nick := r.FormValue("nickname")
				login := r.FormValue("login")
				email := r.FormValue("email")
				pass := r.FormValue("password")
				msg, err := db.AddRecord(table, fname, sname, nick, login, email, pass)
				if err != nil {
					messages.Error(err).Render(r.Context(), w)
					fmt.Printf("Error when try to add user: %v\n	Table: %s\n	Body: %v\n", err, r.Body)
					return
				}
				messages.Message(msg).Render(r.Context(), w)
				return
			case "videofiles":
				videofile := &models.Videofile{}
				json.NewDecoder(r.Body).Decode(videofile)
				msg, err := db.AddRecord(table, videofile.Filename, videofile.Uploader, videofile.Size, videofile.Duration)
				if err != nil {
					fmt.Printf("Error when try to add user: %v\n	Table: %s\n	Body: %v\n", err, table, videofile)
				}
				w.Write([]byte(msg))
				return
			case "broadcasts":
				broadcast := &models.Broadcast{}
				json.NewDecoder(r.Body).Decode(broadcast)
				msg, err := db.AddRecord(table, broadcast.BroadcastStartTime.String(), broadcast.BroadcastEndTime.String())
				if err != nil {
					fmt.Printf("Error when try to add user: %v\n	Table: %s\n	Body: %v\n", err, table, broadcast)
				}
				w.Write([]byte(msg))
				return
			case "broadcasts_files":
				bf := &models.BroadcastFiles{}
				json.NewDecoder(r.Body).Decode(bf)
				msg, err := db.AddRecord(table, bf.BroadcastId, bf.VideofileId)
				if err != nil {
					fmt.Printf("Error when try to add user: %v\n	Table: %s\n	Body: %v\n", err, table, bf)
				}
				w.Write([]byte(msg))
				return
			case "broadcasts_users":

				bu := &models.BroadcastUsers{}
				json.NewDecoder(r.Body).Decode(bu)
				msg, err := db.AddRecord(table, bu.BroadcastId, bu.UserId)
				if err != nil {
					fmt.Printf("Error when try to add user: %v\n	Table: %s\n	Body: %v\n", err, table, bu)
				}
				w.Write([]byte(msg))
				return
			}
			messages.Message(msg).Render(r.Context(), w)
		case "delete":
			msg, err := db.DeleteRecord(table, r.FormValue("id"))
			if err != nil {
				messages.Error(err).Render(r.Context(), w)
			}
			messages.Message(msg).Render(r.Context(), w)
		case "update":
			msg, err := db.UpdateRecord(table, r.FormValue("id"), r.FormValue("fname"), r.FormValue("sname"), r.FormValue("nickname"), r.FormValue("login"), r.FormValue("email"), r.FormValue("password"))
			if err != nil {
				messages.Error(err).Render(r.Context(), w)
			}
			messages.Message(msg).Render(r.Context(), w)
		}
	}
}

func ButtonHandler(w http.ResponseWriter, r *http.Request) {
	button := r.PathValue("Button")
	tp := r.PathValue("Type")
	// log.Println("button: ", button, " type: ", tp)
	switch button {
	case "delete":
		switch tp {
		default:
			messages.Error(fmt.Errorf("table %s not found", tp)).Render(r.Context(), w)
			return
		}
	case "update":
		// buttons.UpdateTableDialog().Render(r.Context(), w)
		return
	default:
		messages.Error(fmt.Errorf("Wtf you want from me????")).Render(r.Context(), w)
	}

}
