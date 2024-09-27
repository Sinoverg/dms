package structs

type User struct {
	Uuid      string
	FirstName string
	LastName  string
	Login     string
	Email     string
	PassHash  string
}
type File struct {
	Uuid     string
	Name     string
	Size     string
	Duration string
}
type Broadcast struct {
	FileUuid       string
	UserUuid       string
	StartQueueTime string
	EndQueueTime   string
	Date           string
}
