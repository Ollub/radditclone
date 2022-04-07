package users

type User struct {
	Id       int `sql:"AUTO_INCREMENT"`
	Name     string
	PassHash string
}

// incoming data for creating users
type UserIn struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}
