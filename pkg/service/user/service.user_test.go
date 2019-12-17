package user

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nandawinata/entry-task/pkg/helper/bcrypt"
	"github.com/nandawinata/entry-task/pkg/service/user/data"
	"github.com/stretchr/testify/mock"
)

const (
	successGetUsername = "Object1"
	failGetUsername    = "Object2"
	anyError           = "ObjAnyError"
	correctPass        = "correctPass"
	failPass           = "failPass"
)

type MyMockedObject struct {
	mock.Mock
}

var mockUser *data.User

func init() {
	hashPass, _ := bcrypt.HashPassword(correctPass)
	mockUser = &data.User{
		Username: successGetUsername,
		Nickname: successGetUsername,
		Password: hashPass,
	}
}

func (m *MyMockedObject) GetUserById(id uint64) (*data.User, error) {
	args := m.Called(id)
	if id == 0 {
		return mockUser, nil
	}
	return nil, args.Error(1)
}

func (m *MyMockedObject) GetUserByUsername(username string) (*data.User, error) {
	m.Called(username)
	switch username {
	case successGetUsername:
		return mockUser, nil
	case failGetUsername:
		return nil, nil
	}

	return nil, errors.New("Any Error")
}

func (m *MyMockedObject) InsertUser(user data.User) (*data.User, error) {
	return mockUser, nil
}

func (m *MyMockedObject) InsertUserBulk(payload data.UserBulkPayload) error {
	args := m.Called(payload)
	return args.Error(1)
}

func (m *MyMockedObject) UpdateNickname(user data.User) error {
	return nil
}

func (m *MyMockedObject) UpdatePhoto(user data.User) error {
	return nil
}

func TestLogin(t *testing.T) {
	// User not found
	testObj := new(MyMockedObject)
	userService := UserService{testObj}

	loginPayload := LoginPayload{
		Username: failGetUsername,
		Password: failPass,
	}

	testObj.On("GetUserByUsername", failGetUsername).Return(nil, nil)
	_, err := userService.Login(loginPayload)

	if err != nil {
		fmt.Println(err.Error())
	}
	testObj.AssertExpectations(t)
	// End User not found

	// Error Get Username
	testObjThree := new(MyMockedObject)
	userService = UserService{testObjThree}

	loginPayload = LoginPayload{
		Username: anyError,
		Password: failPass,
	}

	testObjThree.On("GetUserByUsername", anyError).Return(nil, mock.AnythingOfType("error"))
	_, err = userService.Login(loginPayload)

	if err != nil {
		fmt.Println(err.Error())
	}
	testObjThree.AssertExpectations(t)
	// End Error Get Username

	// IncorrectPassword
	testObjTwo := new(MyMockedObject)
	userService = UserService{testObjTwo}

	loginPayload = LoginPayload{
		Username: successGetUsername,
		Password: failPass,
	}

	testObjTwo.On("GetUserByUsername", successGetUsername).Return(mockUser, nil)
	_, err = userService.Login(loginPayload)

	if err != nil {
		fmt.Println(err.Error())
	}
	testObjTwo.AssertExpectations(t)
	// End IncorrectPassword

	// Success
	testObjOne := new(MyMockedObject)
	userService = UserService{testObjOne}
	loginPayload = LoginPayload{
		Username: successGetUsername,
		Password: correctPass,
	}
	testObjOne.On("GetUserByUsername", successGetUsername).Return(mockUser, nil)
	_, err = userService.Login(loginPayload)

	if err != nil {
		fmt.Println(err.Error())
	}
	testObjOne.AssertExpectations(t)
	// End Success
}
