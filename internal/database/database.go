package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DBController struct {
	db *sql.DB
}

func NewController() (*DBController,error) {
	db, err := connectDatabase()
	return &DBController{
		db: db,
	},err
}

func connectDatabase() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host= %s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5432", os.Getenv("POSTGRES_USER_NAME"), os.Getenv("POSTGRES_USER_PASSWORD"), os.Getenv("DATABASE_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil,err
	}
	return db, nil
}

func (d *DBController) SelectAll(table string) (*sql.Rows, error) {
	rows, err := d.db.Query("SELECT * FROM " + table + ";")
	if err != nil {
		return &sql.Rows{}, err
	}
	return rows, err
}

func (dbc *DBController) AddVideofile(f, u, s, d string)(error){
		result, err := dbc.db.Exec("INSERT INTO videofiles ( filename, uploader, size, duration) VALUES ($1, $2, $3, $4)",
		 f,u,s,d)
		if err != nil {
			return err
		}
		if _, err = result.RowsAffected(); err != nil {
			return err
		}
		return nil
}
func (d *DBController) AddUser(fn, sn, nick, login, email, pass string)(error){
		result, err := d.db.Exec("INSERT INTO users (firstname, secondaryname, nickname, login, email, password) VALUES ($1, $2, $3, $4, $5, $6)",
		 fn,sn,nick,login,email,pass)
		if err != nil {
			return err
		}
		if _, err = result.RowsAffected(); err != nil {
			return err
		}
		return nil
}
func (d *DBController)AddBroadcast(bst, bet time.Time)(error){
		result, err := d.db.Exec("INSERT INTO broadcasts (broadcaststarttime, broadcastendtime) VALUES ($1, $2);", bst,bet)
		if err != nil {
			return err
		}
		if _, err = result.RowsAffected(); err != nil {
			return err
		}
		return nil
}
func (d *DBController)AddBroadcastUser(bID,uID string)(error){
		result, err := d.db.Exec("INSERT INTO broadcast_users (broadcastid, userid) VALUES ($1, $2)", bID,uID)
		if err != nil {
			return err
		}
		if _, err = result.RowsAffected(); err != nil {
			return err
		}
		return nil
}
func (d *DBController)AddBroadcastFile(bID, fID string)(error){
		result, err := d.db.Exec("INSERT INTO broadcast_files (broadcastid, fileid) VALUES ($1, $2)", bID,fID)
		if err != nil {
			return err
		}
		if _, err = result.RowsAffected(); err != nil {
			return err
		}
		return nil
}
func (d *DBController) DeleteRecord(table, id string) (error) {
	result, err := d.db.Exec("DELETE FROM $1 WHERE id = $2", table,id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (d *DBController) DeleteBroadcastFile(bID, id string)(error){
	result, err := d.db.Exec("DELETE FROM broadcasts_files WHERE broadcastid = $1 AND fileid = $2", bID,id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (d *DBController) DeleteBroadcastUser(bID, id string)(error){
	result, err := d.db.Exec("DELETE FROM broadcasts_users WHERE broadcastid = $1 AND userid = $2", bID,id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (d *DBController) UpdateUser(id, fn, sn, nick, login, email, pass string)(error){
		result, err := d.db.Exec("UPDATE users SET firstname= $1,  secondaryname= $2,  nickname= $3, login= $4, email= $5, password= $6,  WHERE userid=$7", 
		fn,sn,nick,login,email,pass,id)
		if err != nil {
			return err
		}
		if _, err = result.RowsAffected(); err != nil {
		return err
		}
		return nil
}
func (db *DBController) UpdateVideoFile(id, f, u, s, d string)(error){
	result, err := db.db.Exec("UPDATE videofiles SET filename= $1,  uploader= $2,  size= $3, duration= $4 WHERE fileid=$5", 
	f,u,s,d,id)
	if err != nil {
			return err
	}
	if _, err = result.RowsAffected();err != nil {
		return err
	}
	return nil
}
func (db *DBController) UpdateBroadcast(id, st, et string)(error){
	result, err := db.db.Exec("UPDATE broadcasts SET broadcaststarttime= $1,  broadcastendtime= $2,  WHERE userid=$3", 
	st,et,id)
	if err != nil {
			return err
	}
if  _,err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}