package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/nandawinata/entry-task/pkg/service/user"
	"github.com/nandawinata/entry-task/pkg/service/user/data"
)

const (
	limit          = 5000000
	limitExecute   = 5000
	thread         = 10
	dummyPath      = "/assets/data/test.txt"
	BaseInsert     = "INSERT IGNORE INTO users(username, nickname, password) VALUES"
	PreparedInsert = "(?,?,?)"
	Delimiter      = ","
)

type UserPool struct {
	Length   int
	UserBulk data.UserBulkPayload
}

var userPool map[int]UserPool
var wg sync.WaitGroup

func init() {
	userPool = make(map[int]UserPool)
}

func main() {
	counter := 0

	for counter < limit {
		poolID := counter % thread

		wg.Add(1)
		go func(poolID int, randomString string) {
			defer wg.Done()
			poolInsertBulk(poolID, randomString)
		}(poolID, strconv.Itoa(counter))
		wg.Wait()

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
