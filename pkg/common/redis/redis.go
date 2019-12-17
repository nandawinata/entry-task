package redis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	eh "github.com/nandawinata/entry-task/pkg/helper/error_handler"
)

var client *redis.Client

type RedisSetting struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Redis struct {
	client *redis.Client
}

func init() {
	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fmt.Printf("In DIRECTORY %s\n", dir)

	jsonFile, err := os.Open(dir + "/configs/redis/redis.json")
	// jsonFile, err := os.Open("/Users/nandaadhiwinata/go/src/github.com/nandawinata/entry-task/configs/redis/redis.json")
	if err != nil {
		panic(err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(err)
	}

	var setting RedisSetting

	err = json.Unmarshal(byteValue, &setting)

	if err != nil {
		panic(err)
	}

	client = redis.NewClient(&redis.Options{
		Addr: setting.Host,
		// Addr:     "127.0.0.1:6379",
		Password: setting.Password,
		DB:       setting.DB,
	})

	if client == nil {
		panic(fmt.Errorf("Redis not set\n"))
	}
}

func New() Redis {
	return Redis{
		client: client,
	}
}

func (r Redis) Set(key string, value interface{}, expired time.Duration) error {
	key = cleanKey(key)
	marshaledValue, err := json.Marshal(value)
	if err != nil {
		return eh.DefaultError(err)
	}

	err = r.client.Set(key, marshaledValue, expired).Err()
	if err != nil {
		return eh.DefaultError(err)
	}

	return nil
}

func (r Redis) Get(key string, value interface{}) error {
	key = cleanKey(key)
	val, err := r.client.Get(key).Result()

	if err == redis.Nil {
		return nil
	} else if err != nil {
		return eh.DefaultError(err)
	}

	res := []byte(val)
	json.Unmarshal(res, value)

	return nil
}

func cleanKey(key string) string {
	key = strings.Replace(strings.ToLower(key), " ", "", -1)
	return key
}
