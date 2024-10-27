package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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
	return db, nil
}

func (d *DBController) SelectAll(table string) (*sql.Rows, error) {
	rows, err := d.db.Query("SELECT * FROM " + table + ";")
	if err != nil {
		log.Println("Error: ", err)
		return &sql.Rows{}, err
	}
	return rows, err
}
func (d *DBController) AddRecord(table string, data ...string) (string, error) {
	switch table {
	case "videofiles":
		result, err := d.db.Exec("INSERT INTO videofiles ( filename, uploader, size, duration) VALUES ($1, $2, $3, $4)", data[0], data[1], data[2], data[3])
		if err != nil {
			return "", err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Sprintln("Rows inserted: ", err.Error()), nil
		}
		return fmt.Sprintln("Rows inserted: ", affected), nil
	case "users":
		result, err := d.db.Exec("INSERT INTO users (fname, sname, nickname, login, email, password) VALUES ($1, $2, $3, $4, $5, $6)", data[0], data[1], data[2], data[3], data[4], data[5])
		if err != nil {
			return "", err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Sprintln("Rows inserted: ", err.Error()), nil
		}
		return fmt.Sprintln("Rows inserted: ", affected), nil
	case "broadcasts":
		log.Printf("INSERT INTO broadcasts (broadcaststarttime, broadcastendtime) VALUES (%s, %s)", data[0], data[1])
		result, err := d.db.Exec("INSERT INTO broadcasts (broadcaststarttime, broadcastendtime) VALUES (?, ?);", data[0], data[1])
		if err != nil {
			return "", err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Sprintln("Rows inserted: ", err.Error()), nil
		}
		return fmt.Sprintln("Rows inserted: ", affected), nil
	case "broadcast_users":
		result, err := d.db.Exec("INSERT INTO broadcast_users (broadcastid, userid) VALUES ($1, $2)", data[0], data[1])
		if err != nil {
			return "", err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Sprintln("Rows inserted: ", err.Error()), nil
		}
		return fmt.Sprintln("Rows inserted: ", affected), nil
	case "broadcast_files":
		result, err := d.db.Exec("INSERT INTO broadcast_files (broadcastid, fileid) VALUES ($1, $2)", data[0], data[1])
		if err != nil {
			return "", err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Sprintln("Rows inserted: ", err.Error()), nil
		}
		return fmt.Sprintln("Rows inserted: ", affected), nil
	}
	return "", fmt.Errorf("table %s not found", table)
}

func (d *DBController) DeleteRecord(table string, data ...string) (string, error) {
	result, err := d.db.Exec("DELETE FROM videofiles WHERE id = $1", data[0])
	if err != nil {
		return "", err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Sprintln("Rows deleted: ", err.Error()), nil
	}
	return fmt.Sprintln("Rows deleted: ", affected), nil
}

func (d *DBController) UpdateRecord(table string, data ...string) (string, error) {
	var result sql.Result
	err := fmt.Errorf("")
	switch table {
	case "users":
		result, err = d.db.Exec("UPDATE users SET firstname= ?,  secondaryname= ?,  nickname= ?, login= ?, email= ?, password= ?,  WHERE userid=?", data[1], data[2], data[3], data[4], data[5], data[6], data[0])
	case "videofiles":
		result, err = d.db.Exec("UPDATE videofiles SET filename= ?,  uploader= ?,  size= ?, duration= ? WHERE fileid=?", data[1], data[2], data[3], data[4], data[0])
	case "broadcasts":
		result, err = d.db.Exec("UPDATE broadcasts SET broadcaststarttime= ?,  broadcastendtime= ?  WHERE userid=?", data[1], data[2], data[0])
	case "broadcasts_users":
		result, err = d.db.Exec("UPDATE broadcasts_users SET broadcastid= ?,  userid= ?  WHERE userid=?", data[1], data[2], data[0])
	case "broadcasts_files":
		result, err = d.db.Exec("UPDATE broadcasts_files SET broadcastid= ?,  fileid= ?  WHERE userid=?", data[1], data[2], data[0])
	}
	if err != nil {
		return "", err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Sprintln("Rows updated: ", err.Error()), nil
	}
	return fmt.Sprintln("Rows updated: ", affected), nil
}
