package utils

import (
	"context"
	"diamond/config"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RDClient redis客户端
var RDClient *redis.Client

func init() {
	host := config.Config.Get("redis.host")
	port := config.Config.Get("redis.port")
	password := config.Config.Get("redis.password").(string)
	dbName := int(config.Config.Get("redis.dbName").(float64))
	poolSize := int(config.Config.Get("redis.poolSize").(float64))

	addr := fmt.Sprintf("%v:%v", host, port)
	RDClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbName,
		PoolSize: poolSize,
	})
}

// SetToken 更新用户token, 默认是24小时过期
func SetToken(uid uint, token string) {
	ctx := context.Background()
	key := fmt.Sprintf("uid_%v_token", uid)
	err := RDClient.Set(ctx, key, token, 24*time.Hour).Err()
	if err != nil {
		panic(err)
	}
}

// GetToken 查询用户token
func GetToken(uid uint) string {
	ctx := context.Background()
	key := fmt.Sprintf("uid_%v_token", uid)
	token, err := RDClient.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return token
}

// DelToken 删除用户token
func DelToken(uid uint) {
	ctx := context.Background()
	key := fmt.Sprintf("uid_%v_token", uid)
	err := RDClient.Del(ctx, key).Err()
	if err != nil {
		panic(err)
	}
}
