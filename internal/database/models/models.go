package models

type Record interface {
}

type Delete struct{
	Id string `json:"id"`
	Table string `json:"table"`
}

type User struct {
	Id       string `json:"id"`
	Fname    string `json:"first_name"`
	Sname    string `json:"second_name"`
	Nickname string `json:"nickname"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Broadcast struct {
	Id                 string    `json:"id"`
		BroadcastStartTime string `json:"b_start_time"`
		BroadcastEndTime   string `json:"b_end_time"`
}

type Videofile struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Uploader string `json:"uploader"`
	Size     string `json:"size"`
	Duration string `json:"duration"`
}

type BroadcastFiles struct {
	BroadcastId string `json:"broadcast_id"`
	VideofileId string `json:"videofile_id"`
}

type BroadcastUsers struct {
	BroadcastId string `json:"broadcast_id"`
	UserId      string `json:"user_id"`
}
