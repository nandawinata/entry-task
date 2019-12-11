package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/nandawinata/entry-task/pkg/handler/constants"
	eh "github.com/nandawinata/entry-task/pkg/helper/error_handler"
	mw "github.com/nandawinata/entry-task/pkg/helper/middleware"
	"github.com/nandawinata/entry-task/pkg/helper/upload_file"
	"github.com/nandawinata/entry-task/pkg/service/user"
)

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	var payload user.RegisterPayload
	err = json.Unmarshal(b, &payload)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	user, err := user.New().Register(payload)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return user, nil
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) (interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	var payload user.LoginPayload
	err = json.Unmarshal(b, &payload)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	response, err := user.New().Login(payload)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return response, nil
}

func UpdateProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params, token *mw.TokenPayload) (interface{}, error) {
	if token == nil {
		return nil, eh.NewError(http.StatusUnauthorized, "Invalid token")
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	var payload user.UpdatePayload
	err = json.Unmarshal(b, &payload)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	payload.ID = token.ID

	err = user.New().Update(payload)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return nil, nil
}

func UpdatePhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, token *mw.TokenPayload) (interface{}, error) {
	if token == nil {
		return nil, eh.NewError(http.StatusUnauthorized, "Invalid token")
	}

	dir, err := os.Getwd()

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	savedPath := dir + constants.SAVED_PATH
	filePath, err := upload_file.UploadFile(w, r, savedPath)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	payload := user.UpdatePayload{
		ID:    token.ID,
		Photo: filePath,
	}

	err = user.New().Update(payload)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return nil, nil
}
