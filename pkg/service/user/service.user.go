package user

import (
	"net/http"

	"github.com/nandawinata/entry-task/pkg/helper/bcrypt"
	eh "github.com/nandawinata/entry-task/pkg/helper/error_handler"
	"github.com/nandawinata/entry-task/pkg/helper/middleware"
	"github.com/nandawinata/entry-task/pkg/helper/upload_file"
	"github.com/nandawinata/entry-task/pkg/service/user/data"
)

type UserService struct {
	data data.UserData
}

func New() UserService {
	return UserService{
		data: data.MysqlUserData{},
	}
}

type RegisterPayload struct {
	Username       string `json:"username"`
	Nickname       string `json:"nickname"`
	Password       string `json:"password"`
	RetypePassword string `json:"retype_password"`
}

func (s UserService) Register(payload RegisterPayload) (*data.UserOutput, error) {
	if payload.Password != payload.RetypePassword {
		return nil, eh.NewError(http.StatusBadRequest, "Password and Retype Password not match")
	}

	user, err := s.GetUserByUsername(payload.Username)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	if user != nil {
		return nil, eh.NewError(http.StatusBadRequest, "Username already taken")
	}

	hashPass, err := bcrypt.HashPassword(payload.Password)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	savedUser := data.User{
		Username: payload.Username,
		Nickname: payload.Nickname,
		Password: hashPass,
	}

	temp, err := s.data.InsertUser(savedUser)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	userOutput := data.UserOutput{
		ID:       temp.ID,
		Username: temp.Username,
		Nickname: temp.Nickname,
		Photo:    temp.Photo,
	}

	return &userOutput, nil
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string          `json:"token"`
	User  data.UserOutput `json:"user"`
}

func (s UserService) Login(payload LoginPayload) (*LoginResponse, error) {
	user, err := s.GetUserByUsername(payload.Username)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	if user == nil {
		return nil, eh.NewError(http.StatusBadRequest, "Username not found")
	}

	if !bcrypt.CheckPasswordHash(user.Password, payload.Password) {
		return nil, eh.NewError(http.StatusBadRequest, "Incorrect password")
	}

	tokenPayload := middleware.TokenPayload{
		ID:       user.ID,
		Username: user.Username,
	}

	token, err := middleware.GenerateJwt(tokenPayload)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	if token == nil {
		return nil, eh.NewError(http.StatusBadRequest, "Failed to generate token")
	}

	loginResponse := LoginResponse{
		Token: *token,
		User: data.UserOutput{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Photo:    user.Photo,
		},
	}

	return &loginResponse, nil
}

func (s UserService) GetUserByUsername(username string) (*data.User, error) {
	user, err := s.data.GetUserByUsername(username)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return user, nil
}

func (s UserService) GetUserById(id uint64) (*data.User, error) {
	user, err := s.data.GetUserById(id)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return user, nil
}

type UpdatePayload struct {
	ID       uint64  `json:"id" db:"id"`
	Nickname *string `json:"nickname"`
	Photo    *string `json:"photo" db:"photo"`
}

func (s UserService) Update(payload UpdatePayload) error {
	user, err := s.GetUserById(payload.ID)

	if err != nil {
		return eh.NewError(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return eh.NewError(http.StatusBadRequest, "User not found")
	}

	updatedUser := data.User{
		ID: payload.ID,
	}

	if payload.Nickname != nil && *payload.Nickname != user.Nickname {
		updatedUser.Nickname = *payload.Nickname
		err := s.data.UpdateNickname(updatedUser)

		if err != nil {
			return eh.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	if payload.Photo != nil {
		updatedUser.Photo = payload.Photo
		err := s.data.UpdatePhoto(updatedUser)

		if err != nil {
			return eh.NewError(http.StatusInternalServerError, err.Error())
		}

		if user.Photo != nil {
			_ = upload_file.DeleteFile(*user.Photo)
		}
	}

	return nil
}
