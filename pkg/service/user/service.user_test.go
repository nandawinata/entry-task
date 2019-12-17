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
	successID             = uint64(0)
	failID                = uint64(1)
	errorID               = uint64(2)
	successGetUsername    = "Object1"
	failGetUsername       = "Object2"
	anyError              = "ObjAnyError"
	correctPass           = "correctPass"
	failPass              = "failPass"
	updateNicknameSuccess = "anyNickname"
	updateNicknameFailed  = "anyNicknameOne"
)

type MyMockedObject struct {
	mock.Mock
}

var (
	mockUser       *data.User
	mockUserInsert *data.User
)

func init() {
	hashPass, _ := bcrypt.HashPassword(correctPass)
	mockUser = &data.User{
		Username: successGetUsername,
		Nickname: successGetUsername,
		Password: hashPass,
	}

	mockUserInsert = &data.User{
		ID:       uint64(1),
		Username: failGetUsername,
		Nickname: failGetUsername,
		Password: hashPass,
	}
}

func (m *MyMockedObject) GetUserById(id uint64) (*data.User, error) {
	m.Called(id)
	switch id {
	case successID:
		return mockUser, nil
	case failID:
		return nil, nil
	}

	return nil, errors.New("Any Error")
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
	m.Called(user)
	if user.Username == failGetUsername {
		return mockUserInsert, nil
	}

	return nil, errors.New("Any Error")
}

func (m *MyMockedObject) InsertUserBulk(payload data.UserBulkPayload) error {
	args := m.Called(payload)
	return args.Error(1)
}

func (m *MyMockedObject) UpdateNickname(user data.User) error {
	m.Called(user)

	if user.Nickname == updateNicknameSuccess {
		return nil
	}

	return errors.New("Any Error")
}

func (m *MyMockedObject) UpdatePhoto(user data.User) error {
	m.Called(user)

	if *user.Photo == updateNicknameSuccess {
		fmt.Println("KAKAKK")
		return nil
	}

	return errors.New("Any Error")
}

func TestRegister(t *testing.T) {
	// Password and Retype not match
	testObj := new(MyMockedObject)
	userService := UserService{testObj}

	registerPayload := RegisterPayload{
		Username:       failGetUsername,
		Password:       failPass,
		RetypePassword: correctPass,
	}

	_, err := userService.Register(registerPayload)

	if err != nil {
		fmt.Println(err.Error())
	}
	testObj.AssertExpectations(t)
	// End Password and Retype not match

	// Error Get Username
	testObjThree := new(MyMockedObject)
	userService = UserService{testObjThree}

	registerPayload = RegisterPayload{
		Username:       anyError,
		Password:       failPass,
		RetypePassword: failPass,
	}

	testObjThree.On("GetUserByUsername", anyError).Return(nil, mock.AnythingOfType("error"))
	_, err = userService.Register(registerPayload)

	if err != nil {
		fmt.Println(err.Error())
	}
	testObjThree.AssertExpectations(t)
	// End Error Get Username
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

func TestUpdate(t *testing.T) {
	// User not found
	testObj := new(MyMockedObject)
	userService := UserService{testObj}

	updatePayload := UpdatePayload{
		ID: failID,
	}

	testObj.On("GetUserById", failID).Return(nil, nil)
	err := userService.Update(updatePayload)

	if err != nil {
		fmt.Println(err.Error())
	}
	testObj.AssertExpectations(t)
	// End User not found

	// Error GetUserByID
	testObjOne := new(MyMockedObject)
	userService = UserService{testObjOne}

	updatePayload = UpdatePayload{
		ID: errorID,
	}

	testObjOne.On("GetUserById", errorID).Return(nil, nil)
	err = userService.Update(updatePayload)

	if err != nil {
		fmt.Println(err.Error())
	}
	testObjOne.AssertExpectations(t)
	// End Error GetUserByID

	// Update nickname success
	testObjThree := new(MyMockedObject)
	userService = UserService{testObjThree}

	nickname := updateNicknameSuccess
	updatePayload = UpdatePayload{
		ID:       successID,
		Nickname: &nickname,
	}

	testObjThree.On("GetUserById", successID).Return(mockUser, nil)

	updateNicknamePayload := data.User{
		ID:       successID,
		Nickname: nickname,
	}
	testObjThree.On("UpdateNickname", updateNicknamePayload).Return(nil)
	err = userService.Update(updatePayload)

	if err != nil {
		fmt.Println(err.Error())
	}
	testObjThree.AssertExpectations(t)
	// End Update nickname success
}
