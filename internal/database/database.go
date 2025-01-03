package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Cirqach/dms/internal/database/models"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type DBController struct {
	DB *sql.DB
}

func NewController() (*DBController, error) {
	db, err := connectDatabase()
	return &DBController{
		DB: db,
	}, err
}

func connectDatabase() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host= %s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("POSTGRES_USER_NAME"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("DATABASE_NAME"))
	time.Sleep(5 * time.Second)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (d *DBController) SelectUser(id string) (*sql.Rows, error) {
	rows, err := d.DB.Query("SELECT * FROM users WHERE userid = $1;", id)
	if err != nil {
		return &sql.Rows{}, err
	}
	return rows, err
}
func (d *DBController) SelectBroadcast(id string) (*sql.Rows, error) {
	rows, err := d.DB.Query("SELECT * FROM broadcasts WHERE broadcastid = $1;", id)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DBController) SelectVideofile(id string) (*sql.Rows, error) {
	rows, err := d.DB.Query("SELECT * FROM videofiles WHERE fileid = $1;", id)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DBController) SelectBroadcastFile(broadcastID, fileID string) (*sql.Rows, error) {
	rows, err := d.DB.Query("SELECT * FROM broadcasts_files WHERE broadcastid = $1 and fileID = $2;", broadcastID, fileID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DBController) SelectBroadcastUser(broadcastID, userID string) (*sql.Rows, error) {
	rows, err := d.DB.Query("SELECT * FROM broadcasts_users WHERE broadcastid = $1 and userID = $2;", broadcastID, userID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DBController) SelectAll(table string) (*sql.Rows, error) {
	rows, err := d.DB.Query("SELECT * FROM " + table + ";")
	if err != nil {
		return &sql.Rows{}, err
	}
	return rows, err
}

func (dbc *DBController) AddVideofile(f string, u int, s, d string) error {
	result, err := dbc.DB.Exec("INSERT INTO videofiles ( filename, uploader, size, duration) VALUES ($1, $2, $3, $4)",
		f, u, s, d)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (d *DBController) AddUser(fn, sn, nick, login, email, pass string) error {
	result, err := d.DB.Exec("INSERT INTO users (firstname, secondaryname, nickname, login, email, password) VALUES ($1, $2, $3, $4, $5, $6)",
		fn, sn, nick, login, email, pass)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (d *DBController) AddBroadcast(bst, bet time.Time) error {
	// Format the time values to "YYYY-MM-DD HH:MM:SS"
	formattedBST := bst.Format("2006-01-02 15:04:05")
	formattedBET := bet.Format("2006-01-02 15:04:05")

	// Execute the SQL query with formatted time strings
	result, err := d.DB.Exec(
		"INSERT INTO broadcasts (broadcaststarttime, broadcastendtime) VALUES ($1, $2);",
		formattedBST, formattedBET,
	)
	if err != nil {
		return err
	}

	// Check if rows were affected
	if _, err = result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (d *DBController) AddBroadcastUser(bID, uID string) error {
	// Convert string IDs to integers
	broadcastID, err := strconv.Atoi(bID)
	if err != nil {
		return fmt.Errorf("invalid broadcast ID: %w", err)
	}

	userID, err := strconv.Atoi(uID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Execute the query with integer values
	result, err := d.DB.Exec("INSERT INTO broadcasts_users (broadcastid, userid) VALUES ($1, $2)", broadcastID, userID)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (d *DBController) AddBroadcastFile(bID, fID int) error {
	result, err := d.DB.Exec("INSERT INTO broadcasts_files (broadcastid, fileid) VALUES ($1, $2)", bID, fID)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (d *DBController) DeleteUsers(id string) error {
	result, err := d.DB.Exec("DELETE FROM users WHERE userid = $1", id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (d *DBController) DeleteBroadcasts(id string) error {
	result, err := d.DB.Exec("DELETE FROM broadcasts WHERE broadcastid = $1", id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	log.Info().Msg(fmt.Sprintf("Count of affected rows: %d", affected))
	return nil
}

func (d *DBController) DeleteVideofiles(id string) error {
	result, err := d.DB.Exec("DELETE FROM videofiles WHERE fileid = $1", id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (d *DBController) DeleteBroadcastFile(bID, fID int) error {
	result, err := d.DB.Exec("DELETE FROM broadcasts_files WHERE broadcastid = $1 AND fileid = $2", bID, fID)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (d *DBController) DeleteBroadcastUser(bID, id int) error {
	result, err := d.DB.Exec("DELETE FROM broadcasts_users WHERE broadcastid = $1 AND userid = $2", bID, id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (d *DBController) UpdateUser(user models.User) error {
	// log.Info().Msg(fmt.Sprintf("Updating user %s %s %s %s %s %s %s",id,fn,sn,nick,login,email,pass))
	result, err := d.DB.Exec("UPDATE users SET firstname=$1,  secondaryname=$2,  nickname=$3, login=$4, email=$5, password=$6  WHERE userid=$7;",
		user.Fname, user.Sname, user.Nickname, user.Login, user.Email, user.Password, user.Id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (db *DBController) UpdateVideoFile(id,u int, f, s, d string) error {
	result, err := db.DB.Exec("UPDATE videofiles SET filename=$1,  uploader=$2,  size=$3, duration=$4 WHERE fileid=$5",
		f, u, s, d, id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (db *DBController) UpdateBroadcast(b models.Broadcast) error {
	result, err := db.DB.Exec("UPDATE broadcasts SET broadcaststarttime=$1,  broadcastendtime=$2 WHERE broadcastid=$3",
		b.StartTime, b.EndTime, b.Id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (db *DBController) UpdateBroadcastUsers(old, new models.BroadcastUsers) error {
	result, err := db.DB.Exec("UPDATE broadcasts_users SET broadcastid=$1,  userid=$2 WHERE broadcastid=$3 and userid=$4",
		new.BroadcastId, new.UserId, old.BroadcastId, old.UserId)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
func (db *DBController) UpdateBroadcastFiles(old, new models.BroadcastFiles) error {
	result, err := db.DB.Exec("UPDATE broadcasts_files SET broadcastid=$1,  fileid=$2 WHERE broadcastid=$3 and fileid=$4",
		new.BroadcastId, new.VideofileId, old.BroadcastId, old.VideofileId)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (db *DBController) SearchVideofiles() {

}
