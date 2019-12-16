package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type DBSetting struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Scheme   string `json:"scheme"`
}

var DB *sqlx.DB

func init() {
	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fmt.Printf("In DIRECTORY %s", dir)
	fmt.Println()

	jsonFile, err := os.Open(dir + "/configs/database/database.json")

	if err != nil {
		panic(err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(err)
	}

	var setting DBSetting

	err = json.Unmarshal(byteValue, &setting)

	if err != nil {
		panic(err)
	}

	DB, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", setting.Username, setting.Password, setting.Host, setting.Scheme))

	if err != nil {
		panic(err)
	}
}

func GetDB() *sqlx.DB {
	return DB
}
