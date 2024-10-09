package handler

import (
	"net/http"

	"github.com/Cirqach/dms/cmd/templ/messages"
	"github.com/Cirqach/dms/cmd/templ/tables"
	"github.com/Cirqach/dms/internal/database"
	"github.com/a-h/templ"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {

}

func TableHandler(db *database.DBController, table string) templ.Component {
	switch table {
	case "files":
		files, err := db.SelectFiles()
		if err != nil {
			return messages.Error(err)
		}
		return tables.Files(files)
	case "users":
		users, err := db.SelectUsers()
		if err != nil {
			return messages.Error(err)
		}
		return tables.Users(users)
	}
	broadcast, err := db.SelectBroadcast()
	if err != nil {
		return messages.Error(err)
	}
	return tables.Broadcast(broadcast)
}
