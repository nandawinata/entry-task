package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	eh "github.com/nandawinata/entry-task/pkg/helper/error_handler"
)

type HandleResponse struct {
	Code    int         `json:"code"`
	Message *string     `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseFormatter func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) (interface{}, error)

func ResponseJson(h ResponseFormatter) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		data, err := h(w, r, p)

		if err == nil {
			response := HandleResponse{
				Code: 200,
				Data: data,
			}
			body, _ := json.Marshal(response)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(body)
			return
		}

		e, ok := err.(*eh.ErrorString)

		if !ok {
			response := HandleResponse{
				Code: 500,
			}
			body, _ := json.Marshal(response)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(body)
			return
		}

		fmt.Println(e.Stacktrace())

		errMessage := e.Error()

		response := HandleResponse{
			Code:    e.Code(),
			Message: &errMessage,
		}

		body, _ := json.Marshal(response)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Code())
		w.Write(body)
		return
	}
}
