package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	eh "github.com/nandawinata/entry-task/pkg/helper/error_handler"
	jsf "github.com/nandawinata/entry-task/pkg/helper/json"
	mw "github.com/nandawinata/entry-task/pkg/helper/middleware"
)

type HandlerWithToken func(writer http.ResponseWriter, request *http.Request, params httprouter.Params, token *mw.TokenPayload) (interface{}, error)

func ValidateJwt(next HandlerWithToken) jsf.ResponseFormatter {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
		token, err := mw.ValidateJwt(r)

		if err != nil {
			return nil, eh.DefaultError(err)
		}

		data, err := next(w, r, p, token)

		return data, err
	}
}
