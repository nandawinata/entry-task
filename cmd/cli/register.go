package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nandawinata/entry-task/pkg/service/user"
	"github.com/nandawinata/entry-task/pkg/service/user/data"
)

const (
	limit          = 5000000
	limitExecute   = 1000
	thread         = 10
	dummyPath      = "/data/test.txt"
	BaseInsert     = "INSERT INTO users(username, nickname, password) VALUES"
	PreparedInsert = "(?,?,?)"
	Delimiter      = ","
)

type UserPool struct {
	Length   int
	UserBulk data.UserBulkPayload
}

var userPool map[int]UserPool

func init() {
	userPool = make(map[int]UserPool)
}

func main() {
	counter := 0

	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fmt.Printf("LOAD FROM [%s]\n", dir+dummyPath)
	inFile, err := os.Open(dir + dummyPath)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		if counter >= 5000000 {
			break
		}

		poolID := counter % thread

		go func(poolID int, randomString string) {
			poolInsertBulk(poolID, randomString)
		}(poolID, scanner.Text())

		counter++
	}
}

func poolInsertBulk(poolID int, randomString string) {
	pool, ok := userPool[poolID]

	if !ok {
		pool = UserPool{
			Length: 0,
			UserBulk: data.UserBulkPayload{
				Query: BaseInsert,
			},
		}
	}

	if pool.Length > 0 {
		pool.UserBulk.Query = pool.UserBulk.Query + Delimiter
	}
	pool.UserBulk.Query = pool.UserBulk.Query + PreparedInsert
	pool.UserBulk.Params = append(pool.UserBulk.Params, randomString, randomString, randomString)
	pool.Length++
	fmt.Printf("Append data to POOL[%d] --> VALUES[%s]\n", poolID, randomString)

	if pool.Length == limitExecute {
		fmt.Printf("Execute BULK INSERT --> POOL[%d]\n", poolID)
		err := user.New().InsertUserBulk(pool.UserBulk)

		if err != nil {
			panic(err)
		}

		pool = UserPool{
			Length: 0,
			UserBulk: data.UserBulkPayload{
				Query: BaseInsert,
			},
		}
	}

	userPool[poolID] = pool
}
