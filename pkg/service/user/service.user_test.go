package user

import (
	"fmt"
	"testing"

	"github.com/nandawinata/entry-task/pkg/service/user/data"
	"github.com/stretchr/testify/mock"
)

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) GetUserById(id uint64) (*data.User, error) {
	args := m.Called(id)
	result := data.User{}
	return &result, args.Error(1)
}

func (m *MyMockedObject) GetUserByUsername(username string) (*data.User, error) {
	args := m.Called(username)
	result := data.User{}
	return &result, args.Error(1)
}

func (m *MyMockedObject) InsertUser(user data.User) (*data.User, error) {
	args := m.Called(user)
	result := data.User{}
	return &result, args.Error(1)
}

func (m *MyMockedObject) InsertUserBulk(payload data.UserBulkPayload) error {
	args := m.Called(payload)
	return args.Error(1)
}

func (m *MyMockedObject) UpdateNickname(user data.User) error {
	args := m.Called(user)
	return args.Error(1)
}

func (m *MyMockedObject) UpdatePhoto(user data.User) error {
	args := m.Called(user)
	return args.Error(1)
}

func TestRegister(t *testing.T) {
	testObj := new(MyMockedObject)

	testObj.On("GetUserById", 123).Return(nil, nil)

	userService := UserService{testObj}
	a, err := userService.Register(RegisterPayload{})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(a)

	testObj.AssertExpectations(t)
}
