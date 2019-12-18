package user

import (
	"fmt"
	"net/http"

	"github.com/nandawinata/entry-task/pkg/common/redis"
	"github.com/nandawinata/entry-task/pkg/helper/bcrypt"
	eh "github.com/nandawinata/entry-task/pkg/helper/error_handler"
	"github.com/nandawinata/entry-task/pkg/helper/middleware"
	"github.com/nandawinata/entry-task/pkg/helper/upload_file"
	"github.com/nandawinata/entry-task/pkg/service/user/constants"
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
	var user *data.User

	redisService := redis.New()
	keyUsername := fmt.Sprintf(constants.KEY_USERNAME, username)
	err := redisService.Get(keyUsername, &user)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	if user != nil {
		return user, nil
	}

	user, err = s.data.GetUserByUsername(username)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	err = redisService.Set(keyUsername, user, constants.REDIS_EXPIRE)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	if user != nil {
		keyID := fmt.Sprintf(constants.KEY_USER_ID, int(user.ID))
		err = redisService.Set(keyID, user, constants.REDIS_EXPIRE)

		if err != nil {
			return nil, eh.DefaultError(err)
		}
	}

	return user, nil
}

func (s UserService) GetUserById(id uint64, refreshCache bool) (*data.User, error) {
	var user *data.User

	redisService := redis.New()
	keyID := fmt.Sprintf(constants.KEY_USER_ID, int(id))

	if !refreshCache {
		err := redisService.Get(keyID, &user)

		if err != nil {
			return nil, eh.DefaultError(err)
		}

		if user != nil {
			return user, nil
		}
	}

	user, err := s.data.GetUserById(id)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	err = redisService.Set(keyID, user, constants.REDIS_EXPIRE)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	if user != nil {
		keyUsername := fmt.Sprintf(constants.KEY_USERNAME, user.Username)
		err = redisService.Set(keyUsername, user, constants.REDIS_EXPIRE)

		if err != nil {
			return nil, eh.DefaultError(err)
		}
	}

	return user, nil
}

type UpdatePayload struct {
	ID       uint64  `json:"id" db:"id"`
	Nickname *string `json:"nickname"`
	Photo    *string `json:"photo" db:"photo"`
}

func (s UserService) Update(payload UpdatePayload) error {
	user, err := s.GetUserById(payload.ID, false)

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

	s.GetUserById(payload.ID, true)

	return nil
}

func (s UserService) UserToUserOutput(payload data.User) data.UserOutput {
	userOutput := data.UserOutput{
		ID:       payload.ID,
		Username: payload.Username,
		Nickname: payload.Nickname,
		Photo:    payload.Photo,
	}

	return userOutput
}

func (s UserService) InsertUserBulk(payload data.UserBulkPayload) error {
	return s.data.InsertUserBulk(payload)
}
