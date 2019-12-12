package database

import (
	"fmt"

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
	// dir, err := os.Getwd()

	// if err != nil {
	// 	panic(err)
	// }

	// jsonFile, err := os.Open(dir + "/configs/database/database.json")

	// if err != nil {
	// 	panic(err)
	// }

	// byteValue, err := ioutil.ReadAll(jsonFile)

	// if err != nil {
	// 	panic(err)
	// }

	// var setting DBSetting

	// err = json.Unmarshal(byteValue, &setting)

	// if err != nil {
	// 	panic(err)
	// }

	DB, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", "Angka1234", "localhost:3306", "entry_task"))
	DB.SetMaxIdleConns(8)

	if err != nil {
		panic(err)
	}
}

func GetDB() *sqlx.DB {
	return DB
}
