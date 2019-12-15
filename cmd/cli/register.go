package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/nandawinata/entry-task/pkg/helper/bcrypt"
	"github.com/nandawinata/entry-task/pkg/service/user"
	"github.com/nandawinata/entry-task/pkg/service/user/data"
)

type UserPool struct {
	Length   int
	UserBulk data.UserBulkPayload
}

const (
	BaseInsert     = "INSERT IGNORE INTO users(username, nickname, password) VALUES"
	PreparedInsert = "(?,?,?)"
	Delimiter      = ","
)

var (
	userPool     map[int]UserPool
	limitExecute int
	thread       int
	wg           sync.WaitGroup
)

func init() {
	userPool = make(map[int]UserPool)
}

func main() {
	counter := 0

	limit, _ := strconv.Atoi(os.Args[1])
	limitExecute, _ = strconv.Atoi(os.Args[2])
	thread, _ = strconv.Atoi(os.Args[3])

	if limitExecute <= 0 || thread <= 0 {
		panic(fmt.Errorf("Arguments not valid"))
	}

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

	finalExecute()
}

func poolInsertBulk(poolID int, randomString string) {
	pool, ok := userPool[poolID]

	if !ok {
		pool = resetPool(poolID)
	}

	if pool.Length > 0 {
		pool.UserBulk.Query = pool.UserBulk.Query + Delimiter
	}
	pool.UserBulk.Query = pool.UserBulk.Query + PreparedInsert
	password, _ := bcrypt.HashPassword(randomString)
	pool.UserBulk.Params = append(pool.UserBulk.Params, randomString, randomString, password)
	pool.Length++
	fmt.Printf("Append data to POOL[%d] --> VALUES[%s]\n", poolID, randomString)

	if pool.Length < limitExecute {
		userPool[poolID] = pool
		return
	}

	err := executePool(poolID)
	if err != nil {
		panic(err)
	}
	resetPool(poolID)
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

func resetPool(poolID int) UserPool {
	userPool[poolID] = UserPool{
		Length: 0,
		UserBulk: data.UserBulkPayload{
			Query: BaseInsert,
		},
	}

	return userPool[poolID]
}

func finalExecute() {
	for poolID := 0; poolID < thread; poolID++ {
		executePool(poolID)
		resetPool(poolID)
	}
}
