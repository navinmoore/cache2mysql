package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdsFailoverClient *redis.Client
	err               error
)

func NewFailoverClient() *redis.Client {
	rdsFailoverClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "master",
		SentinelAddrs: []string{"x.x.x.x:26379", "xx.xx.xx.xx:26379", "xxx.xxx.xxx.xxx:26379"},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdsFailoverClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("redis.NewFailoverClient error=", err)
	}
	return rdsFailoverClient
}

func FailoverClientInstance() *redis.Client {
	if rdsFailoverClient != nil {
		return rdsFailoverClient
	}
	rdsLock.Lock()
	defer rdsLock.Unlock()
	if rdsFailoverClient != nil {
		return rdsFailoverClient
	}
	return NewFailoverClient()
}
