package upload_file

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	eh "github.com/nandawinata/entry-task/pkg/helper/error_handler"
)

var allowedType map[string]bool

func init() {
	allowedType = make(map[string]bool)
	allowedType["image/jpeg"] = true
	allowedType["image/jpg"] = true
	allowedType["image/png"] = true
}

func UploadFile(w http.ResponseWriter, r *http.Request, savedPath string) (*string, error) {
	if r.Method != http.MethodPost {
		return nil, eh.NewError(http.StatusBadRequest, "Invalid method request")
	}

	file, handle, err := r.FormFile("file")

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	defer file.Close()

	mimeType := handle.Header.Get("Content-Type")

	_, ok := allowedType[mimeType]

	if !ok {
		return nil, eh.NewError(http.StatusBadRequest, "Invalid file type")
	}

	ext := filepath.Ext(handle.Filename)
	handle.Filename = strconv.Itoa(int(time.Now().Unix())) + ext

	filePath, err := saveFile(w, file, handle, savedPath)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return filePath, nil
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader, savedPath string) (*string, error) {
	uploadedFile, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	if _, err := os.Stat(savedPath); os.IsNotExist(err) {
		_ = os.Mkdir(savedPath, 0700)
	}

	filePath := savedPath + handle.Filename

	err = ioutil.WriteFile(filePath, uploadedFile, 0666)

	if err != nil {
		return nil, eh.DefaultError(err)
	}

	return &filePath, nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)

	if err != nil {
		return eh.DefaultError(err)
	}

	return nil
}
