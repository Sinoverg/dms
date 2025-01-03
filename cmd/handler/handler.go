package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Cirqach/dms/cmd/templ/tables"
	"github.com/Cirqach/dms/cmd/word"
	"github.com/Cirqach/dms/internal/database"
	"github.com/Cirqach/dms/internal/database/models"
	"github.com/lukasjarosch/go-docx"
	"github.com/rs/zerolog/log"
)


func TableHandler(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		table := r.PathValue("table")
		rows, err := db.SelectAll(table)
		if err != nil {
			log.Err(err).Msg("Error selecting all records from " + table + " table")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal error, try later"))
			return
		}
		switch table {
		case "videofiles":
			files := make([]models.Videofile, 0)
			for rows.Next() {
				videofile := models.Videofile{}
				err := rows.Scan(&videofile.Id, &videofile.Filename, &videofile.Uploader, &videofile.Size, &videofile.Duration)
				if err != nil {
					log.Err(err).Msg("Error scanning videofiles rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				files = append(files, videofile)
			}
			users :=make([]models.User, 0)
			userRows, err := db.SelectAll("users")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from users table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}
			for userRows.Next() {
				user := models.User{}
				err := userRows.Scan(&user.Id, &user.Fname, &user.Sname, &user.Nickname, &user.Login, &user.Email, &user.Password)
				if err != nil {
					log.Err(err).Msg("Error scanning users rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				users = append(users, user)
			}
			tables.Files(files,users).Render(r.Context(), w)
			return
		case "users":
			users := make([]models.User, 0)
			for rows.Next() {
				user := models.User{}
				err := rows.Scan(&user.Id, &user.Fname, &user.Sname, &user.Nickname, &user.Login, &user.Email, &user.Password)
				if err != nil {
					log.Err(err).Msg("Error scanning users rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				users = append(users, user)
			}
			tables.Users(users).Render(r.Context(), w)
			return
		case "broadcasts":
			broadcasts := make([]models.Broadcast, 0)
			for rows.Next() {
				broadcast := models.Broadcast{}
				err := rows.Scan(&broadcast.Id, &broadcast.StartTime, &broadcast.EndTime)
				if err != nil {
					log.Err(err).Msg("Error scanning broadcasts rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}

				broadcasts = append(broadcasts, broadcast)
			}

			tables.Broadcast(broadcasts).Render(r.Context(), w)
			return
		case "broadcasts_users":
			broadcastsUsers := make([]models.BroadcastUsers, 0)
			for rows.Next() {
				bu := models.BroadcastUsers{}
				err := rows.Scan(&bu.BroadcastId, &bu.UserId)
				if err != nil {
					log.Err(err).Msg("Error scanning broadcasts_users rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				broadcastsUsers = append(broadcastsUsers, bu)
			}
			broadcasts := make([]models.Broadcast, 0)
			broadcastRow, err := db.SelectAll("broadcasts")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from broadcasts table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}
			for broadcastRow.Next() {
				b := models.Broadcast{}
				err := broadcastRow.Scan(&b.Id, &b.StartTime, &b.EndTime)
				if err != nil {
					log.Err(err).Msg("Error scanning broadcasts rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				
				broadcasts = append(broadcasts, b)
			}
			users := make([]models.User, 0)
			userRow, err := db.SelectAll("users")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from users table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}
			for userRow.Next() {
				user := models.User{}
				err := userRow.Scan(&user.Id, &user.Fname, &user.Sname, &user.Nickname, &user.Login, &user.Email, &user.Password)
				if err != nil {
					log.Err(err).Msg("Error scanning users rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				users = append(users, user)
			}
			tables.Broadcasts_users(broadcastsUsers, users, broadcasts).Render(r.Context(), w)
			return
		case "broadcasts_files":
			broadcastFiles := make([]models.BroadcastFiles, 0)
			for rows.Next() {
				bf := models.BroadcastFiles{}
				err := rows.Scan(&bf.BroadcastId, &bf.VideofileId)
				if err != nil {
					log.Err(err).Msg("Error scanning broadcasts_files rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				broadcastFiles = append(broadcastFiles, bf)
			}
			broadcasts := make([]models.Broadcast, 0)
			broadcastRow, err := db.SelectAll("broadcasts")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from broadcasts table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}

			for broadcastRow.Next() {
				b := models.Broadcast{}
				err := broadcastRow.Scan(&b.Id, &b.StartTime, &b.EndTime)
				if err != nil {
					log.Err(err).Msg("Error scanning broadcasts rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				
				broadcasts = append(broadcasts, b)
			}
			videofiles := make([]models.Videofile, 0)
			videofileRow, err := db.SelectAll("videofiles")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from videofiles table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}

			for videofileRow.Next() {
				f := models.Videofile{}
				err := videofileRow.Scan(&f.Id, &f.Filename, &f.Uploader, &f.Size, &f.Duration)
				if err != nil {
					log.Err(err).Msg("Error scanning videofiles rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				videofiles = append(videofiles, f)
			}
			tables.Broadcasts_files(broadcastFiles, broadcasts, videofiles).Render(r.Context(), w)
			return
		}
	}
}

func ApiAddUser(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			log.Err(err).Msg("Error decoding user from request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err := db.AddUser(user.Fname, user.Sname, user.Nickname, user.Login, user.Email, user.Password); err != nil {
			if strings.Contains(err.Error(),"value violates unique constraint \"users_email_key\""){
			log.Err(err).Msg("user with email already exists")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("User with this email already exists"))
			return
			}
			if strings.Contains(err.Error(),"violates unique constraint \"users_login_key\""){
			log.Err(err).Msg("user with this login already exists")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("User with this login already exists"))
			return
			}
			log.Err(err).Msg("Error adding user to database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		p := &docx.PlaceholderMap{
			"do": "registration",
			"doed": "register new user",
			"data": fmt.Sprintf("New user %s %s registered with login %s, email %s and password: %s",user.Fname, user.Sname,user.Login, user.Email, user.Password),
		}
		filename, err := word.Generate(*p)
		if err != nil {
			log.Err(err).Msg("Error generating docx file")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal error, try later"))
			return
		}
		log.Info().Msg(fmt.Sprintf("User added: %v", user))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User added successfully. You can download register cheque here: <a href=/static/docx/" + filename + "><strong>Download<strong></a>"))
	}
}
func ApiAddVideofile(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		videofile := &models.Videofile{}
		if err := json.NewDecoder(r.Body).Decode(videofile); err != nil {
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
		log.Info().Msg(fmt.Sprintf("Videofile added: %v", videofile))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Videofile added successfully"))
	}
}
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

	
		if b.EndTime.After(b.StartTime) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Broadcast start time must be before end time"))
			return
		}

		// Call AddBroadcast with parsed times
		if err := db.AddBroadcast(b.StartTime, b.EndTime); err != nil {
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
func ApiAddBroadcastsUsers(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse and validate input values
		broadcastID := r.FormValue("broadcastId")
		userID := r.FormValue("userId")
		log.Debug().Msg(fmt.Sprintf("Adding broadcasts_users record: broadcastId=%s, userId=%s", broadcastID, userID))
		// Check for missing values
		if broadcastID == "" || userID == "" {
			http.Error(w, "Missing broadcastId or userId", http.StatusBadRequest)
			return
		}

		// Attempt to add the record
		if err := db.AddBroadcastUser(broadcastID, userID); err != nil {
			log.Err(err).Msg("Error adding broadcasts_users to database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		log.Info().Msg(fmt.Sprintf("Broadcasts_users record added: broadcastId=%s, userId=%s", broadcastID, userID))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Broadcasts_users record added successfully"))
	}
}

func ApiAddBroadcastsFiles(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		broadcastID, err := strconv.Atoi(r.FormValue("broadcastId"))
		if err != nil {
			log.Err(err).Msg("Error converting broadcastID to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		fileid, err := strconv.Atoi(r.FormValue("fileId"))
		if err != nil {
			log.Err(err).Msg("Error converting fileid to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err := db.AddBroadcastFile(broadcastID, fileid); err != nil {
			log.Err(err).Msg("Error adding broadcasts_files to database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Info().Msg(fmt.Sprintf("Broadcasts_files record added: broadcastId=%d, fileId=%d", broadcastID, fileid))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Broadcasts_files record added successfully"))

	}
}
func ApiDelete(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		table := r.PathValue("table")
		id := r.PathValue("id")
		switch table {
		case "users":
			if err := db.DeleteUsers(id); err != nil {
				log.Err(err).Msg("Error deleting user from database")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			log.Info().Msg(fmt.Sprintf("User %v was deleted", id))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User record was successfully deleted"))
		case "videofiles":
			if err := db.DeleteVideofiles(id); err != nil {
				log.Err(err).Msg("Error deleting videofile from database")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			log.Info().Msg(fmt.Sprintf("Videofile %v was deleted", id))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Videofile record was successfully deleted"))
		case "broadcasts":
			if err := db.DeleteBroadcasts(id); err != nil {
				log.Err(err).Msg("Error deleting broadcast from database")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			log.Info().Msg(fmt.Sprintf("Broadcast %v was deleted", id))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Broadcast record was successfully deleted"))
		case "broadcasts_files":
			videofileid, err := strconv.Atoi(r.URL.Query().Get("videofile_id"))
			if err != nil {
				log.Err(err).Msg("Error converting videofileid to int")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			broadcastID, err := strconv.Atoi(id)
			if err != nil {
				log.Err(err).Msg("Error converting broadcastID to int")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			if err := db.DeleteBroadcastFile(broadcastID, videofileid); err != nil {
				log.Err(err).Msg("Error deleting broadcasts_files record from database")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			log.Info().Msg(fmt.Sprintf("Broadcasts_files %v was deleted", id))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Broadcasts_files record was successfully deleted"))
		case "broadcasts_users":
			uid := r.URL.Query().Get("user_id")
			userid, err := strconv.Atoi(uid)
			if err != nil {
				log.Err(err).Msg("Error converting userid to int")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			broadcastID, err := strconv.Atoi(id)
			if err != nil {
				log.Err(err).Msg("Error converting broadcastID to int")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			if err := db.DeleteBroadcastUser(broadcastID, userid); err != nil {
				log.Err(err).Msg("Error deleting broadcasts_users record from database")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			log.Info().Msg(fmt.Sprintf("Broadcasts_users %v was deleted", id))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Broadcasts_users record was successfully deleted"))

		}
	}
}
func ApiDeleteBroadcastsUsers(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bu := &models.BroadcastUsers{}
		if err := json.NewDecoder(r.Body).Decode(bu); err != nil {
			log.Err(err).Msg("Error decoding broadcasts_users from request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err := db.DeleteBroadcastUser(bu.BroadcastId, bu.UserId); err != nil {
			log.Err(err).Msg("Error deleting broadcasts_users from database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Info().Msg(fmt.Sprintf("Broadcasts_users %v was deleted", bu))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Broadcasts_users record was successfully deleted"))
	}
}
func ApiDeleteBroadcastsFiles(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bid := r.PathValue("id")
		fid, err := strconv.Atoi(r.URL.Query().Get("videofile_id"))
		if err != nil {
			log.Err(err).Msg("Error converting fileid to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		broadcastID, err := strconv.Atoi(bid)
		if err != nil {
			log.Err(err).Msg("Error converting broadcastID to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err := db.DeleteBroadcastFile(broadcastID, fid); err != nil {
			log.Err(err).Msg("Error deleting broadcasts_files record from database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Info().Msg(fmt.Sprintf("Broadcasts_files %v was deleted", bid))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Broadcasts_files record was successfully deleted"))
	}
}

func ApiUpdateBroadcastsFiles(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		oldbID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Err(err).Msg("Error converting oldbID to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		oldfid, err := strconv.Atoi(r.URL.Query().Get("videofile_id"))
		if err != nil {
			log.Err(err).Msg("Error converting oldfid to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		oldBF := models.BroadcastFiles{
			BroadcastId: oldbID,
			VideofileId: oldfid,
		}
		newbID, err := strconv.Atoi(r.FormValue("broadcastId"))
		if err != nil {
			log.Err(err).Msg("Error converting newbID to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		newfid, err := strconv.Atoi(r.FormValue("videofileId"))
		if err != nil {
			log.Err(err).Msg("Error converting newfid to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		newBF := models.BroadcastFiles{
			BroadcastId: newbID,
			VideofileId: newfid,
		}
		if err := db.UpdateBroadcastFiles(oldBF, newBF); err != nil {
			log.Err(err).Msg("Error updating broadcasts_files in database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Info().Msg(fmt.Sprintf("Broadcasts_files record updated to %v", newBF))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Broadcasts_files was successfully updated"))

	}
}
func ApiUpdateUser(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Err(err).Msg("Error converting id to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		u := models.User{
			Id:       id,
			Fname:    r.FormValue("first_name"),
			Sname:    r.FormValue("second_name"),
			Nickname: r.FormValue("nickname"),
			Login:    r.FormValue("login"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		if err := db.UpdateUser(u); err != nil {
			log.Err(err).Msg("Error updating user in database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Info().Msg(fmt.Sprintf("User record updated to %v", u))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User was successfully updated"))
	}
}
func ApiUpdateVideofile(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Err(err).Msg("Error converting id to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		uploader,err := strconv.Atoi(r.FormValue("uploader"))

if err != nil {
			log.Err(err).Msg("Error converting uploader to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		v := models.Videofile{
			Id:       id,
			Filename: r.FormValue("filename"),
			Uploader: uploader,
			Size:     r.FormValue("size"),
			Duration: r.FormValue("duration"),
		}
		if err := db.UpdateVideoFile(v.Id, v.Uploader, v.Filename, v.Size, v.Duration); err != nil {
			log.Err(err).Msg("Error updating videofile in database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Info().Msg(fmt.Sprintf("Videofile record updated to %v", v))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Videofile was successfully updated"))
	}
}
func ApiUpdateBroadcast(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Err(err).Msg("Error converting id to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		starttime,err := time.Parse(time.DateTime, r.FormValue("broadcast_start_time"))
		if err != nil {
			log.Err(err).Msg("Error parsing start time")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		endtime,err := time.Parse(time.DateTime, r.FormValue("broadcast_end_time"))
		if err != nil {
			log.Err(err).Msg("Error parsing end time")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		b := models.Broadcast{
			Id:        id,
			StartTime: starttime,
			EndTime:   endtime,
		}
		if err := db.UpdateBroadcast(b); err != nil {
			log.Err(err).Msg("Error updating broadcast in database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Info().Msg(fmt.Sprintf("Broadcast record updated to %v", b))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Broadcast was successfully update"))
	}
}

func ApiUpdateBroadcastUsers(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bid, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Err(err).Msg("Error converting id to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		uid := r.URL.Query().Get("user_id")
		userid, err := strconv.Atoi(uid)
		if err != nil {
			log.Err(err).Msg("Error converting userid to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		newbid, err := strconv.Atoi(r.FormValue("broadcastId"))
		if err != nil {
			log.Err(err).Msg("Error converting broadcastID to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		newuid, err := strconv.Atoi(r.FormValue("userId"))
		if err != nil {
			log.Err(err).Msg("Error converting userid to int")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		oldBU := models.BroadcastUsers{
			BroadcastId: bid,
			UserId:      userid,
		}
		bu := models.BroadcastUsers{
			BroadcastId: newbid,
			UserId:      newuid,
		}

		if err := db.UpdateBroadcastUsers(oldBU, bu); err != nil {
			log.Err(err).Msg("Error updating broadcasts_users in database")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Info().Msg(fmt.Sprintf("Broadcasts_users record updated to %v", bu))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Broadcasts_users was successfully updated"))

	}
}

func Search(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		table := r.PathValue("table")
		query := r.FormValue("search")
		log.Debug().Msg(fmt.Sprintf("Searching for '%s' in table '%s'", query, table))
		var rows *sql.Rows
		var err error
		switch table {
		case "users":
			rows, err = db.DB.Query(`SELECT * FROM users WHERE nickname ILIKE $1 OR 
			login ILIKE $1  OR 
			email ILIKE $1 OR  
			firstname ILIKE $1 OR 
			secondaryname ILIKE $1 OR 
			nickname ILIKE $1 OR 
			secondaryname ILIKE $1`, "%"+query+"%")
		case "videofiles":
			rows, err = db.DB.Query("SELECT * FROM videofiles WHERE filename ILIKE $1", "%"+query+"%")
		case "broadcasts":
			rows, err = db.DB.Query(`SELECT * FROM broadcasts 
WHERE CAST(broadcaststarttime AS TEXT) ILIKE $1 
   OR CAST(broadcastendtime AS TEXT) ILIKE $1
`, "%"+query+"%")
		default:
			http.Error(w, "Invalid table name", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		switch table {
		case "users":
			users := []models.User{}
			for rows.Next() {
				var u models.User
				if err := rows.Scan(&u.Id, &u.Fname, &u.Sname, &u.Nickname, &u.Login, &u.Email, &u.Password); err != nil {
					http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
					return
				}
				users = append(users, u)
			}
			log.Debug().Msg(fmt.Sprintf("Users founded by query(%s): %v", query, users))
			tables.UsersTable(users).Render(r.Context(), w)
		case "videofiles":
			videofiles := []models.Videofile{}
			for rows.Next() {
				var v models.Videofile
				if err := rows.Scan(&v.Id, &v.Filename, &v.Uploader, &v.Size, &v.Duration); err != nil {
					http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
					return
				}
				videofiles = append(videofiles, v)
			}
			log.Debug().Msg(fmt.Sprintf("Videofiles founded by query(%s): %v", query, videofiles))
			users := make([]models.User, 0)
			userRows, err := db.SelectAll("users")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from users table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}
			for userRows.Next() {
				user := models.User{}
				err := userRows.Scan(&user.Id, &user.Fname, &user.Sname, &user.Nickname, &user.Login, &user.Email, &user.Password)
				if err != nil {
					log.Err(err).Msg("Error scanning users rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				users = append(users, user)
			}

			tables.VideofilesTable(videofiles,users).Render(r.Context(), w)
		case "broadcasts":
			broadcasts := []models.Broadcast{}
			for rows.Next() {
				var b models.Broadcast
				if err := rows.Scan(&b.Id, &b.StartTime, &b.EndTime); err != nil {
					http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
					return
				}
				broadcasts = append(broadcasts, b)
			}
			log.Debug().Msg(fmt.Sprintf("Broadcasts founded by query(%s): %v", query, broadcasts))
			tables.BroadcastsTable(broadcasts).Render(r.Context(), w)
		}
	}
}
func Filter(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		table := r.PathValue("table")
		// i pass query values not through url, but through form
		// make for every table they own filter logic
		// my tables in database here:
		/*
		CREATE TABLE Users (
    UserId SERIAL PRIMARY KEY,
    FirstName VARCHAR(50) NOT NULL,
    SecondaryName VARCHAR(50),
    Nickname VARCHAR(50),
    Login VARCHAR(50) UNIQUE NOT NULL,
    Email VARCHAR(100) UNIQUE NOT NULL,
    Password VARCHAR(255) NOT NULL
);

CREATE TABLE VideoFiles (
    FileId SERIAL PRIMARY KEY,
    Filename VARCHAR(255) NOT NULL,
    Uploader INTEGER REFERENCES Users(UserId) ON DELETE CASCADE, -- Cascade delete if the uploader is deleted
    Size BIGINT NOT NULL,
    Duration INTEGER NOT NULL
);

CREATE TABLE Broadcasts (
    BroadcastId SERIAL PRIMARY KEY,
    BroadcastStartTime TIMESTAMP NOT NULL,
    BroadcastEndTime TIMESTAMP NOT NULL
);

CREATE TABLE Broadcasts_files (
    BroadcastId INTEGER REFERENCES Broadcasts(BroadcastId) ON DELETE CASCADE, -- Cascade delete if a broadcast is deleted
    FileId INTEGER REFERENCES VideoFiles(FileId) ON DELETE CASCADE,          -- Cascade delete if a video file is deleted
    PRIMARY KEY (BroadcastId, FileId)
);

CREATE TABLE Broadcasts_users (
    BroadcastId INTEGER REFERENCES Broadcasts(BroadcastId) ON DELETE CASCADE, -- Cascade delete if a broadcast is deleted
    UserId INTEGER REFERENCES Users(UserId) ON DELETE CASCADE,               -- Cascade delete if a user is deleted
    PRIMARY KEY (BroadcastId, UserId)
);
		*/
		switch table {
			case "users":
				fname := r.FormValue("fname")
				sname := r.FormValue("sname")
				nickname := r.FormValue("nickname")
				login := r.FormValue("login")
				email := r.FormValue("email")
				password := r.FormValue("password")
				query := "SELECT * FROM users WHERE 1=1"
				if fname != "" {
					query += fmt.Sprintf(" AND firstname ILIKE '%s'", fname)

				}
				if sname != "" {
					query += fmt.Sprintf(" AND secondaryname ILIKE '%s'", sname)
				}
				if nickname != "" {
					query += fmt.Sprintf(" AND nickname ILIKE '%s'", nickname)
				}
				if login != "" {
					query += fmt.Sprintf(" AND login ILIKE '%s'", login)
				}
				if email != "" {
					query += fmt.Sprintf(" AND email ILIKE '%s'", email)
				}
				if password != "" {
					query += fmt.Sprintf(" AND password ILIKE '%s'", password)
				}
				rows, err := db.DB.Query(query)
				if err != nil {
					http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
					return
				}
				defer rows.Close()
				users := []models.User{}
				for rows.Next() {
					var u models.User
					if err := rows.Scan(&u.Id, &u.Fname, &u.Sname, &u.Nickname, &u.Login, &u.Email, &u.Password); err != nil {
						http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
						return
					}
					users = append(users, u)
				}
				tables.UsersTable(users).Render(r.Context(), w)
			case "videofiles":
	filename := r.FormValue("filename")
	uploader := r.FormValue("uploader")
	bottomSize := r.FormValue("size_bottom")
	topSize := r.FormValue("size_top")
	bottomDuration := r.FormValue("duration_bottom")
	topDuration := r.FormValue("duration_top")


	query := "SELECT * FROM videofiles WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	// Append conditions for text fields using ILIKE
	if filename != "" {
		query += fmt.Sprintf(" AND filename ILIKE $%d", argIdx)
		args = append(args, "%"+filename+"%")
		argIdx++
	}

	// Append conditions for numeric fields
	if uploader != "" {
		query += fmt.Sprintf(" AND uploader = $%d", argIdx)
		args = append(args, uploader) // Note: Ensure uploader is a valid number
		argIdx++
	}
	if bottomSize != "" {
		query += fmt.Sprintf(" AND size >= $%d", argIdx)
		size, err := strconv.ParseInt(bottomSize, 10, 64)
		if err != nil {
			http.Error(w, "Invalid size value", http.StatusBadRequest)
			return
		}
		args = append(args, size)
		argIdx++
	}
	if topSize != "" {
		query += fmt.Sprintf(" AND size <= $%d", argIdx)
		size, err := strconv.ParseInt(topSize, 10, 64)
		if err != nil {
			http.Error(w, "Invalid size value", http.StatusBadRequest)
			return
		}
		args = append(args, size)
		argIdx++
	}

if bottomDuration != "" {
		query += fmt.Sprintf(" AND duration >= $%d", argIdx)
		duration, err := strconv.ParseInt(bottomDuration, 10, 64)
		if err != nil {
			http.Error(w, "Invalid duration value", http.StatusBadRequest)
			return
		}
		args = append(args, duration)
		argIdx++
	}
	if topDuration != "" {
		query += fmt.Sprintf(" AND duration <= $%d", argIdx)
		duration, err := strconv.ParseInt(topDuration, 10, 64)
		if err != nil {
			http.Error(w, "Invalid duration value", http.StatusBadRequest)
			return
		}
		args = append(args, duration)
		argIdx++
	}


	log.Debug().Msg(fmt.Sprintf("Tryna filter the videofiles by this query: %s, args: %v", query, args))

	// Execute the query
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Process the results
	videofiles := []models.Videofile{}
	for rows.Next() {
		var v models.Videofile
		if err := rows.Scan(&v.Id, &v.Filename, &v.Uploader, &v.Size, &v.Duration); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		videofiles = append(videofiles, v)
	}

	log.Debug().Msg(fmt.Sprintf("Videofiles founded by query: %v", videofiles))

	// Retrieve users for rendering
	users := make([]models.User, 0)
	userRows, err := db.SelectAll("users")
	if err != nil {
		log.Err(err).Msg("Error selecting all records from users table")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal error, try later"))
		return
	}
	defer userRows.Close()
	for userRows.Next() {
		user := models.User{}
		if err := userRows.Scan(&user.Id, &user.Fname, &user.Sname, &user.Nickname, &user.Login, &user.Email, &user.Password); err != nil {
			log.Err(err).Msg("Error scanning users rows")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal error, try later"))
			return
		}
		users = append(users, user)
	}

	tables.VideofilesTable(videofiles, users).Render(r.Context(), w)

			case "broadcasts":
				starttime := r.FormValue("start_time")
				endtime := r.FormValue("end_time")
				log.Debug().Msg(fmt.Sprintf("Filtering broadcasts: starttime='%s', endtime='%s'", starttime, endtime))
				query := "SELECT * FROM broadcasts WHERE 1=1"
				if starttime != "" {
					query += fmt.Sprintf(" AND BroadcastStartTime >= '%s'", starttime)
				}
				if endtime != "" {
					query += fmt.Sprintf(" AND broadcastendtime <= '%s'", endtime)
				}
				rows, err := db.DB.Query(query)
				if err != nil {
					http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
					return
				}
				defer rows.Close()
				broadcasts := []models.Broadcast{}
				for rows.Next() {
					var b models.Broadcast
					if err := rows.Scan(&b.Id, &b.StartTime, &b.EndTime); err != nil {
						http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
						return
					}
					broadcasts = append(broadcasts, b)
				}
				log.Debug().Msg(fmt.Sprintf("Broadcasts founded by query(from '%s' to '%s'): %v", starttime,endtime, broadcasts))
				tables.BroadcastsTable(broadcasts).Render(r.Context(), w)
			case "broadcasts_users":
				broadcastID := r.FormValue("broadcastId")
				userID := r.FormValue("userId")
				query := "SELECT * FROM broadcasts_users WHERE 1=1"
				if broadcastID != "" {
					query += fmt.Sprintf(" AND broadcastid = %s", broadcastID)
				}
				if userID != "" {
					query += fmt.Sprintf(" AND userid = %s", userID)
				}
				log.Debug().Msg(fmt.Sprintf("Broadcasts_users search by query: %v", query))
				rows, err := db.DB.Query(query)
				if err != nil {
					http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
					return
				}
				defer rows.Close()
				broadcastsUsers := []models.BroadcastUsers{}
				for rows.Next() {
					var bu models.BroadcastUsers
					if err := rows.Scan(&bu.BroadcastId, &bu.UserId); err !=nil {
						http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
						return
					}
					broadcastsUsers = append(broadcastsUsers, bu)
				}
				broadcasts := make([]models.Broadcast, 0)
				broadcastRow, err := db.SelectAll("broadcasts")
				if err != nil {
					log.Err(err).Msg("Error selecting all records from broadcasts table")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				for broadcastRow.Next() {
					b := models.Broadcast{}
					err := broadcastRow.Scan(&b.Id, &b.StartTime, &b.EndTime)
					if err != nil {
						log.Err(err).Msg("Error scanning broadcasts rows")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("Internal error, try later"))
						return
					}
					broadcasts = append(broadcasts, b)
				}
				users := make([]models.User, 0)
				userRow, err := db.SelectAll("users")
				if err != nil {
					log.Err(err).Msg("Error selecting all records from users table")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				for userRow.Next() {
					user := models.User{}
					err := userRow.Scan(&user.Id, &user.Fname, &user.Sname, &user.Nickname, &user.Login, &user.Email, &user.Password)
					if err != nil {
						log.Err(err).Msg("Error scanning users rows")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("Internal error, try later"))
						return
					}
					users = append(users, user)
				}
				log.Debug().Msg(fmt.Sprintf("Broadcasts_users founded by query: %v", broadcastsUsers))
				tables.BroadcastsUsersTable(broadcastsUsers, users, broadcasts).Render(r.Context(), w)
			case "broadcasts_files":
				broadcastID := r.FormValue("broadcastId")
				videofileID := r.FormValue("videofileId")
				query := "SELECT * FROM broadcasts_files WHERE 1=1"
				if broadcastID != "" {
					query += fmt.Sprintf(" AND broadcastid = %s", broadcastID)
				}
				if videofileID != "" {
					query += fmt.Sprintf(" AND fileid = %s", videofileID)
				}
				log.Debug().Msg(fmt.Sprintf("Broadcasts_files search by query: %v", query))
				rows, err := db.DB.Query(query)
				if err != nil {
					http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
					return
				}
				defer rows.Close()
				broadcastFiles := []models.BroadcastFiles{}
				for rows.Next() {
					var bf models.BroadcastFiles
					if err := rows.Scan(&bf.BroadcastId, &bf.VideofileId); err != nil {
						http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
						return
					}
					broadcastFiles = append(broadcastFiles, bf)
				}
				broadcasts := make([]models.Broadcast, 0)
				broadcastRow, err := db.SelectAll("broadcasts")
				if err != nil {
					log.Err(err).Msg("Error selecting all records from broadcasts table")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				for broadcastRow.Next() {
					b := models.Broadcast{}
					err := broadcastRow.Scan(&b.Id, &b.StartTime, &b.EndTime)
					if err != nil {
						log.Err(err).Msg("Error scanning broadcasts rows")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("Internal error, try later"))
						return
					}
					broadcasts = append(broadcasts, b)
				}
				videofiles := make([]models.Videofile, 0)
				videofileRow, err := db.SelectAll("videofiles")
				if err != nil {
					log.Err(err).Msg("Error selecting all records from videofiles table")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internalerror, try later"))
					return
				}
				for videofileRow.Next() {
					f := models.Videofile{}
					err := videofileRow.Scan(&f.Id, &f.Filename, &f.Uploader, &f.Size, &f.Duration)
					if err != nil {
						log.Err(err).Msg("Error scanning videofiles rows")
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("Internal error, try later"))
						return
					}
					videofiles = append(videofiles, f)
				}
				log.Debug().Msg(fmt.Sprintf("Broadcasts_files founded by query: %v", broadcastFiles))
				tables.BroadcastsFilesTable(broadcastFiles, broadcasts, videofiles).Render(r.Context(), w)
		}
	}

}
func Sort(db *database.DBController) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		table := r.PathValue("table")
		sortType := r.PathValue("type")
		sortOrder := r.PathValue("order")
		query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s %s", table, sortType, sortOrder)
		rows, err := db.DB.Query(query)
		if err != nil {
			http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		switch table {
		case "users":
			users := []models.User{}
			for rows.Next() {
				var u models.User
				if err := rows.Scan(&u.Id, &u.Fname, &u.Sname, &u.Nickname, &u.Login, &u.Email, &u.Password); err != nil {
					http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
					return
				}
				users = append(users, u)
			}
			tables.UsersTable(users).Render(r.Context(), w)
		case "videofiles":
			videofiles := []models.Videofile{}
			for rows.Next() {
				var v models.Videofile
				if err := rows.Scan(&v.Id, &v.Filename, &v.Uploader, &v.Size, &v.Duration); err != nil {
					http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
					return
				}
				videofiles = append(videofiles, v)
			}
			users := make([]models.User, 0)
			userRows, err := db.SelectAll("users")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from users table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}
			for userRows.Next() {
				user := models.User{}
				err := userRows.Scan(&user.Id, &user.Fname, &user.Sname, &user.Nickname, &user.Login, &user.Email, &user.Password)
				if err != nil {					log.Err(err).Msg("Error scanning users rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				users = append(users, user)
			}
			tables.VideofilesTable(videofiles,users).Render(r.Context(), w)
		case "broadcasts":
			broadcasts := []models.Broadcast{}
			for rows.Next() {
				var b models.Broadcast
				if err := rows.Scan(&b.Id, &b.StartTime, &b.EndTime); err != nil {
					http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
					return
				}
				broadcasts = append(broadcasts, b)
			}
			tables.BroadcastsTable(broadcasts).Render(r.Context(), w)
		case "broadcasts_users":
			broadcastsUsers := []models.BroadcastUsers{}
			for rows.Next() {
				var bu models.BroadcastUsers
				if err := rows.Scan(&bu.BroadcastId, &bu.UserId); err != nil {
					http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
					return
				}
				broadcastsUsers = append(broadcastsUsers, bu)
			}
			broadcasts := make([]models.Broadcast, 0)
			broadcastRow, err := db.SelectAll("broadcasts")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from broadcasts table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}
			for broadcastRow.Next() {
				b := models.Broadcast{}
				err := broadcastRow.Scan(&b.Id, &b.StartTime, &b.EndTime)
				if err != nil {
					log.Err(err).Msg("Error scanning broadcasts rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				broadcasts = append(broadcasts, b)
			}
			users := make([]models.User, 0)
			userRow, err := db.SelectAll("users")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from users table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}
			for userRow.Next() {
				user := models.User{}
				err := userRow.Scan(&user.Id, &user.Fname, &user.Sname, &user.Nickname, &user.Login, &user.Email, &user.Password)
				if err != nil {
					log.Err(err).Msg("Error scanning users rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				users = append(users, user)
			}
			tables.BroadcastsUsersTable(broadcastsUsers, users, broadcasts).Render(r.Context(), w)
		case "broadcasts_files":
			broadcastFiles := []models.BroadcastFiles{}
			for rows.Next() {
				var bf models.BroadcastFiles
				if err := rows.Scan(&bf.BroadcastId, &bf.VideofileId); err != nil {
					http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
					return
				}
				broadcastFiles = append(broadcastFiles, bf)
			}
			broadcasts := make([]models.Broadcast, 0)
			broadcastRow, err := db.SelectAll("broadcasts")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from broadcasts table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}
			for broadcastRow.Next() {
				b := models.Broadcast{}
				err := broadcastRow.Scan(&b.Id, &b.StartTime, &b.EndTime)
				if err != nil {
					log.Err(err).Msg("Error scanning broadcasts rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				broadcasts = append(broadcasts, b)
			}
			videofiles := make([]models.Videofile, 0)
			videofileRow, err := db.SelectAll("videofiles")
			if err != nil {
				log.Err(err).Msg("Error selecting all records from videofiles table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal error, try later"))
				return
			}
			for videofileRow.Next() {
				f := models.Videofile{}
				err := videofileRow.Scan(&f.Id, &f.Filename, &f.Uploader, &f.Size, &f.Duration)
				if err != nil {
					log.Err(err).Msg("Error scanning videofiles rows")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal error, try later"))
					return
				}
				videofiles = append(videofiles, f)
			}
			tables.BroadcastsFilesTable(broadcastFiles, broadcasts, videofiles).Render(r.Context(), w)
		}

	}
}
