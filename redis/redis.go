package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb     *redis.Client
	rdsLock sync.Mutex
)

// 集群
func NewClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		PoolSize: 100,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("redis.Newclient error=", err)
	}
	return rdb
}

func RDBInstance() *redis.Client {
	if rdb != nil {
		return rdb
	}
	rdsLock.Lock()
	defer rdsLock.Unlock()
	if rdb != nil {
		return rdb
	}
	return NewClient()
}
