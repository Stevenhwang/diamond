package cache

import (
	"context"
	"diamond/misc"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var Cache *redis.Client

func init() {
	host := misc.Config.Get("redis.host")
	port := misc.Config.Get("redis.port")
	password := misc.Config.GetString("redis.password")
	dbName := misc.Config.GetInt("redis.dbName")
	poolSize := misc.Config.GetInt("redis.poolSize")

	addr := fmt.Sprintf("%v:%v", host, port)
	Cache = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbName,
		PoolSize: poolSize,
	})
	if res := Cache.Ping(ctx); res.Err() != nil {
		panic(res.Err())
	}
}

// ban ip for 7 days
func Ban(ip string) error {
	res := Cache.Set(ctx, ip, 1, 168*time.Hour)
	return res.Err()
}

// get baned ip, true means banned ip exists
func GetBan(ip string) (bool, error) {
	_, err := Cache.Get(ctx, ip).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
