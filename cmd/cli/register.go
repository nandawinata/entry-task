package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

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
	query        string
	userPool     map[int]UserPool
	limitExecute int
	pools        int
)

func init() {
	initPool()
}

func main() {
	counter := 0

	limit, _ := strconv.Atoi(os.Args[1])
	limitExecute, _ = strconv.Atoi(os.Args[2])
	pools, _ = strconv.Atoi(os.Args[3])
	staticPass := os.Args[4]

	if limit <= 0 || limitExecute <= 0 || pools <= 0 {
		panic(fmt.Errorf("Arguments not valid"))
	}

	fmt.Printf("LIMIT[%d] | LIMIT EXECUTE [%d] | POOL [%d]\n", limit, limitExecute, pools)
	staticPass, _ = bcrypt.HashPassword(staticPass)

	start := time.Now()
	fmt.Printf("Start at: %s\n", start.String())

	createQuery()
	for counter < limit {
		poolID := counter % pools
		poolInsertBulk(poolID, strconv.Itoa(counter), staticPass)
		counter++
	}
	finalExecute()

	finish := time.Now()
	elapsed := finish.Sub(start).Seconds()

	fmt.Printf("Finish at: %s\n", finish.String())
	fmt.Printf("Elapsed Time: [%f seconds]\n", elapsed)
}

func createQuery() string {
	query = BaseInsert + PreparedInsert

	for i := 1; i < limitExecute; i++ {
		query = query + Delimiter + PreparedInsert
	}

	return query
}

func poolInsertBulk(poolID int, randomString, staticPass string) {
	pool, _ := userPool[poolID]
	pool.UserBulk.Params = append(pool.UserBulk.Params, randomString, randomString, staticPass)
	pool.Length++
	userPool[poolID] = pool

	if pool.Length < limitExecute {
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

	pool.UserBulk.Query = query
	if ok && pool.Length > 0 {
		return user.New().InsertUserBulk(pool.UserBulk)
	}

	return nil
}

func initPool() {
	userPool = make(map[int]UserPool)

	for poolID := 0; poolID < pools; poolID++ {
		resetPool(poolID)
	}
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
	for poolID := 0; poolID < pools; poolID++ {
		executePool(poolID)
	}
}
