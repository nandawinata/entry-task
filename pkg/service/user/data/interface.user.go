package data

type User struct {
	ID       uint64  `json:"id" db:"id"`
	Username string  `json:"username" db:"username"`
	Nickname string  `json:"nickname" db:"nickname"`
	Password string  `json:"password" db:"password"`
	Photo    *string `json:"photo" db:"photo"`
}

type UserOutput struct {
	ID       uint64  `json:"id" db:"id"`
	Username string  `json:"username" db:"username"`
	Nickname string  `json:"nickname" db:"nickname"`
	Photo    *string `json:"photo" db:"photo"`
}

type UserData interface {
	GetUserById(id uint64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	InsertUser(user User) (*User, error)
	UpdateNickname(user User) error
	UpdatePhoto(user User) error
}
