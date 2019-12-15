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
	limitExecute   = 10000
	thread         = 10
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
		if counter == limit-1 {
			finalExecute()
			break
		}
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
		resetPool(poolID)
	}

	if pool.Length > 0 {
		pool.UserBulk.Query = pool.UserBulk.Query + Delimiter
	}
	pool.UserBulk.Query = pool.UserBulk.Query + PreparedInsert
	pool.UserBulk.Params = append(pool.UserBulk.Params, randomString, randomString, randomString)
	pool.Length++
	fmt.Printf("Append data to POOL[%d] --> VALUES[%s]\n", poolID, randomString)

	if pool.Length == limitExecute {
		executePool(poolID)
		resetPool(poolID)
		return
	}

	userPool[poolID] = pool
}

func executePool(poolID int) error {
	pool, ok := userPool[poolID]

	if ok {
		fmt.Printf("Execute BULK INSERT --> POOL[%d]\n", poolID)
		return user.New().InsertUserBulk(pool.UserBulk)
	}

	fmt.Println("POOL not found\n")
	return nil
}

func resetPool(poolID int) {
	userPool[poolID] = UserPool{
		Length: 0,
		UserBulk: data.UserBulkPayload{
			Query: BaseInsert,
		},
	}
}

func finalExecute() {
	for poolID := 0; poolID < thread; poolID++ {
		executePool(poolID)
		resetPool(poolID)
	}
}
