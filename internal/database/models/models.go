package models

import "time"

type Record interface {
}

type Delete struct {
	Id    string `json:"id"`
	Table string `json:"table"`
}

type User struct {
	Id       int    `json:"id"`
	Fname    string `json:"first_name"`
	Sname    string `json:"second_name"`
	Nickname string `json:"nickname"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Broadcast struct {
	Id        int    `json:"id"`
	StartTime time.Time `json:"b_start_time"`
	EndTime   time.Time `json:"b_end_time"`
}

type Videofile struct {
	Id       int    `json:"id"`
	Filename string `json:"filename"`
	Uploader int `json:"uploader"`
	Size     string `json:"size"`
	Duration string `json:"duration"`
}

type BroadcastFiles struct {
	BroadcastId int `json:"broadcast_id"`
	VideofileId int `json:"videofile_id"`
}

type BroadcastUsers struct {
	BroadcastId int `json:"broadcast_id"`
	UserId      int `json:"user_id"`
}
