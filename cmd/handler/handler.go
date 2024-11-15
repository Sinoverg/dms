package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Cirqach/dms/cmd/templ/tables"
	"github.com/Cirqach/dms/internal/database"
	"github.com/Cirqach/dms/internal/database/models"
	"github.com/rs/zerolog/log"
)



func TableHandler(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		table := r.PathValue("table")
		rows, err := db.SelectAll(table)
		if err != nil {
		log.Err(err).Msg("Error selecting all records from " + table +" table")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
		}
		switch table {
		case "videofiles":
			files := make([]models.Videofile, 0)
			for rows.Next() {
				var id, filename, uploader, size, duration string
				err := rows.Scan(&id, &filename, &uploader, &size, &duration)
				if err != nil {
					log.Err(err).Msg("Error scanning videofiles rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
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
					log.Err(err).Msg("Error scanning users rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
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
				var id, starttime, endtime string 
				err := rows.Scan(&id, &starttime, &endtime)
				if err != nil {
					log.Err(err).Msg("Error scanning broadcasts rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
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
					log.Err(err).Msg("Error scanning broadcasts_users rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
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
					log.Err(err).Msg("Error scanning broadcasts_files rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
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

func ApiAddUser(db *database.DBController) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
				user := &models.User{}
				if err := json.NewDecoder(r.Body).Decode(user); err != nil {
					log.Err(err).Msg("Error decoding user from request body")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				if err := db.AddUser(user.Fname, user.Sname, user.Nickname, user.Login, user.Email, user.Password); err != nil {
					log.Err(err).Msg("Error adding user to database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				log.Info().Msg(fmt.Sprintf("User added: %v",user))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("User added successfully"))
}
}
func ApiAddVideofile(db *database.DBController) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
				videofile := &models.Videofile{}
				if err := json.NewDecoder(r.Body).Decode(videofile); err != nil{
					log.Err(err).Msg("Error decoding videofile from request body")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				if err := db.AddVideofile(videofile.Filename, videofile.Uploader, videofile.Size, videofile.Duration); err != nil {
					log.Err(err).Msg("Error adding videofile to database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				log.Info().Msg(fmt.Sprintf("Videofile added: %v",videofile))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Videofile added successfully"))
	}}
	func ApiAddBroadcast(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			b := &models.Broadcast{}
	
			// Decode JSON body into `b`
			if err := json.NewDecoder(r.Body).Decode(b); err != nil {
				log.Err(err).Msg("Error decoding broadcast from request body")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
	
			// Define your custom format for parsing
			const timeFormat = "2006-01-02T15:04" // Matches the format sent from HTML form

			bst, err := time.Parse(timeFormat, b.BroadcastStartTime)
			if err != nil {
				log.Err(err).Msg("Invalid BroadcastStartTime format")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid BroadcastStartTime format: " + err.Error()))
				return
			}
			
			bet, err := time.Parse(timeFormat, b.BroadcastEndTime)
			if err != nil {
				log.Err(err).Msg("Invalid BroadcastEndTime format")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid BroadcastEndTime format: " + err.Error()))
				return
			}
			if bst.After(bet) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Broadcast start time must be before end time"))
				return
			}
							
			// Call AddBroadcast with parsed times
			if err := db.AddBroadcast(bst, bet); err != nil {
				log.Err(err).Msg("Error adding broadcast to database")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
	
			log.Print(fmt.Sprintf("Broadcast added: %v", b))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Broadcast added successfully"))
		}
	}
func ApiAddBroadcastsUsers(db *database.DBController) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
				bu := &models.BroadcastUsers{}
				if err := json.NewDecoder(r.Body).Decode(bu); err != nil {
					log.Err(err).Msg("Error decoding broadcasts_users from request body")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				if err := db.AddBroadcastUser(bu.BroadcastId, bu.UserId);err != nil {
					log.Err(err).Msg("Error adding broadcasts_users to database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				log.Info().Msg(fmt.Sprintf("Broadcasts_users record added: %v",bu))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Broadcasts_users record added successfully"))
	}}
func ApiAddBroadcastsFiles(db *database.DBController) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
				bf := &models.BroadcastFiles{}
				if err := json.NewDecoder(r.Body).Decode(bf); err != nil{
					log.Err(err).Msg("Error decoding broadcasts_files from request body")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				if err := db.AddBroadcastFile(bf.BroadcastId, bf.VideofileId); err != nil {
					log.Err(err).Msg("Error adding broadcasts_files to database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				log.Info().Msg(fmt.Sprintf("Broadcasts_files record added: %v",bf))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Broadcasts_files record added successfully"))
	}}
func ApiDelete(db *database.DBController) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		switch r.PathValue("table"){
		case "broadcasts_users":
				bu := &models.BroadcastUsers{}
			if err := json.NewDecoder(r.Body).Decode(bu); err != nil {
					log.Err(err).Msg("Error decoding broadcasts_users from request body")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
			}
			if err := db.DeleteBroadcastUser(bu.BroadcastId,bu.UserId); err != nil {
					log.Err(err).Msg("Error deleting broadcasts_users from database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
			}
			log.Info().Msg(fmt.Sprintf("Broadcasts_users %v was deleted",bu))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Broadcasts_users record was successfully deleted"))
		case "broadcasts_files":
				bf := &models.BroadcastFiles{}
				if err := json.NewDecoder(r.Body).Decode(bf); err != nil{
					log.Err(err).Msg("Error decoding broadcasts_files from request body")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
			if err := db.DeleteBroadcastFile(bf.BroadcastId,bf.VideofileId); err != nil {
					log.Err(err).Msg("Error deleting broadcasts_files from database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
			}
			log.Info().Msg(fmt.Sprintf("Broadcasts_files %v was deleted",bf))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Broadcasts_files record was successfully deleted"))
		default:
			table := r.PathValue("table")
			id := r.PathValue("id")
			if err := db.DeleteRecord(table, id); err != nil {
					log.Err(err).Msg("Error deleting record from database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
			log.Info().Msg(fmt.Sprintf("Record with id %d was deleted from %s table",table,id))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(table +" with id "+ id +" was successfully deleted"))
		}
}
}
func ApiUpdateUser(db *database.DBController) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
				u := &models.User{}
				if err := json.NewDecoder(r.Body).Decode(u); err != nil{
					log.Err(err).Msg("Error decoding user from request body")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
			if err := db.UpdateUser(u.Id,u.Fname,u.Sname,u.Nickname,u.Login,u.Email,u.Password); err != nil {
					log.Err(err).Msg("Error updating user in database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
			}
			log.Info().Msg(fmt.Sprintf("User record updated to %v",u))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User was successfully updated"))
}
}
func ApiUpdateVideofile(db *database.DBController) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
				v := &models.Videofile{}
			if err := json.NewDecoder(r.Body).Decode(v); err != nil{
					log.Err(err).Msg("Error decoding videofile from request body")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
			}
			if err := db.UpdateVideoFile(v.Id,v.Filename,v.Uploader,v.Size,v.Duration); err != nil {
					log.Err(err).Msg("Error updating videofile in database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
			}
			log.Info().Msg(fmt.Sprintf("Videofile record updated to %v",v))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Videofile was successfully updated"))
}
}
func ApiUpdateBroadcast(db *database.DBController) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
				b := &models.Broadcast{}
				if err := json.NewDecoder(r.Body).Decode(b); err != nil{
					log.Err(err).Msg("Error decoding broadcast from request body")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
			if err := db.UpdateBroadcast(b.Id,b.BroadcastStartTime,b.BroadcastEndTime); err != nil {
					log.Err(err).Msg("Error updating broadcast in database")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
			}
			log.Info().Msg(fmt.Sprintf("Broadcast record updated to %v",b))
			w.WriteHeader(http.StatusOK)
		w.Write([]byte("Broadcast was successfully update"))
}
}