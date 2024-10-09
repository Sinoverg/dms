package models

import "time"

type User struct {
	Uuid     string
	Fname    string
	Sname    string
	Login    string
	Email    string
	Password string
}

type Broadcast struct {
	Uuid               string
	FileUuid           string
	UserUuid           string
	BroadcastStartTime time.Time
	BroadcastEndTime   time.Time
}

type File struct {
	Uuid     string
	Name     string
	Size     string
	Duration string
}
