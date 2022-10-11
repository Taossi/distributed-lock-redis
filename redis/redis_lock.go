package redis

import (
	"time"

	goredis "github.com/go-redis/redis"
)

const (
	tryTimes int = 3
)

type RedisLock struct {
	key    string
	uuid   string        // uuid作为锁唯一标识
	expire time.Duration // key过期时间
}

func NewRedisLock(key, uuid string, expire time.Duration) *RedisLock {
	return &RedisLock{
		key:    key,
		uuid:   uuid,
		expire: expire,
	}
}

// lock
func (l *RedisLock) Lock() (bool, error) {
	rcli := GetRedisClient()
	return rcli.SetNX(l.key, l.uuid, l.expire).Result()
}

// retry lock
func (l *RedisLock) TryLock() (bool, error) {
	var err error
	var locked bool
	for i := 0; i < tryTimes; i++ {
		locked, err = l.Lock()
		if err == nil && locked {
			return true, nil
		}
		time.Sleep(time.Millisecond * 20)
	}
	return false, err
}

// lua脚本 保证get+del操作原子性
var delScript = goredis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end
`)

// unlock
func (l *RedisLock) UnLock() (bool, error) {
	rcli := GetRedisClient()
	result, err := delScript.Run(rcli, []string{l.key}, l.uuid).Int64()
	return result != 0, err
}

var refreshScript = goredis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("expire", KEYS[1], ARGV[2])
else
	return 0
end
`)

func (l *RedisLock) RefreshLock() (bool, error) {
	rcli := GetRedisClient()
	result, err := refreshScript.Run(rcli, []string{l.key}, l.uuid, int64(l.expire/time.Second)).Int64()
	return result != 0, err
}
