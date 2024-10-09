package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Cirqach/dms/internal/database/models"
	_ "github.com/lib/pq"
)

type DBController struct {
	db *sql.DB
}

func NewController() *DBController {
	db, err := connectDatabase()
	if err != nil {
		log.Fatal("Error making connection with database: ", err)
	}
	return &DBController{
		db: db,
	}
}

func connectDatabase() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host= %s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5432", os.Getenv("POSTGRES_USER_NAME"), os.Getenv("POSTGRES_USER_PASSWORD"), os.Getenv("DATABASE_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db, nil
}

func (database *DBController) SelectFiles() ([]models.File, error) {
	rows, err := database.db.Query("SELECT * FROM files")
	if err != nil {
		return nil, err
	}
	files := make([]models.File, 0)
	for rows.Next() {
		var uuid, name, size, duration string
		err := rows.Scan(uuid, name, size, duration)
		if err != nil {
			return nil, err
		}
		file := models.File{
			Uuid:     uuid,
			Name:     name,
			Size:     size,
			Duration: duration,
		}
		files = append(files, file)
	}
	return files, nil
}
func (database *DBController) SelectUsers() ([]models.User, error) {
	rows, err := database.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	for rows.Next() {
		var uuid, fname, sname, login, email, password string
		err := rows.Scan(uuid, fname, sname, login, email, password)
		if err != nil {
			return nil, err
		}
		user := models.User{
			Uuid:     uuid,
			Fname:    fname,
			Sname:    sname,
			Login:    login,
			Email:    email,
			Password: password,
		}
		users = append(users, user)
	}
	return users, nil
}
func (database *DBController) SelectBroadcast() ([]models.Broadcast, error) {
	rows, err := database.db.Query("SELECT * FROM broadcast")
	if err != nil {
		return nil, err
	}
	broadcast := make([]models.Broadcast, 0)
	for rows.Next() {
		var uuid, fileuuid, useruuid string
		var starttime, endtime time.Time

		err := rows.Scan(uuid, fileuuid, useruuid, starttime, endtime)
		if err != nil {
			return nil, err
		}
		b := models.Broadcast{
			Uuid:               uuid,
			FileUuid:           fileuuid,
			UserUuid:           useruuid,
			BroadcastStartTime: starttime,
			BroadcastEndTime:   endtime,
		}
		broadcast = append(broadcast, b)
	}
	return broadcast, nil
}
