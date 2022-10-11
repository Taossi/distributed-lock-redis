package redis

import (
	"log"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
)

// redis init
func CreateClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		log.Panicf("Redis init failed. err: %s", err)
	}
	return
}

// get redis client
func GetRedisClient() *redis.Client {
	return client
}

func SetHello(key string) error {
	err := client.Set(key, "hello", 0).Err()
	if err != nil {
		return err
	}
	log.Println("set hello success")
	return nil
}

func GetHello(key string) (string, error) {
	value, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func DelHello(key string) error {
	err := client.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}

