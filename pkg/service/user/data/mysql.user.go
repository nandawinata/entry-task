package data

import (
	"database/sql"
	"net/http"

	"github.com/nandawinata/entry-task/pkg/common/database"
	eh "github.com/nandawinata/entry-task/pkg/helper/error_handler"
)

type MysqlUserData struct {
}

const (
	GetUserById       = "SELECT id, username, nickname, password, photo FROM user where id = ?"
	GetUserByUsername = "SELECT id, username, nickname, password, photo FROM user where username = ?"
	InsertUser        = "INSERT INTO USER(username, nickname, password) VALUES(?,?,?)"
	UpdateNickname    = "UPDATE user SET nickname = ? WHERE id = ?"
	UpdatePhoto       = "UPDATE user SET photo = ? WHERE id = ?"
)

func (m MysqlUserData) GetUserById(id uint64) (*User, error) {
	var user User

	stmt, err := database.GetDB().Preparex(GetUserById)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	err = stmt.Get(&user, id)
	defer stmt.Close()

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return &user, nil
}

func (m MysqlUserData) GetUserByUsername(username string) (*User, error) {
	var user User

	stmt, err := database.GetDB().Preparex(GetUserByUsername)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	err = stmt.Get(&user, username)
	defer stmt.Close()

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return &user, nil
}

func (m MysqlUserData) InsertUser(user User) (*User, error) {
	stmt, err := database.GetDB().Preparex(InsertUser)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	result, err := stmt.Exec(user.Username, user.Nickname, user.Password)
	defer stmt.Close()

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	lastInsertedId, err := result.LastInsertId()

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	user.ID = uint64(lastInsertedId)

	return &user, nil
}

func (m MysqlUserData) UpdateNickname(user User) error {
	stmt, err := database.GetDB().Preparex(UpdateNickname)

	if err != nil {
		return eh.NewError(http.StatusInternalServerError, err.Error())
	}

	_, err = stmt.Exec(user.Nickname, user.ID)
	defer stmt.Close()

	if err != nil {
		return eh.NewError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (m MysqlUserData) UpdatePhoto(user User) error {
	stmt, err := database.GetDB().Preparex(UpdatePhoto)

	if err != nil {
		return eh.NewError(http.StatusInternalServerError, err.Error())
	}

	_, err = stmt.Exec(user.Photo, user.ID)

	if err != nil {
		return eh.NewError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
