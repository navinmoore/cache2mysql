package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdsClusterClient *redis.ClusterClient
)

// 集群
func NewRDBClusterClient() *redis.ClusterClient {
	rdsClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
		Password: "",
		PoolSize: 100,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdsClusterClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("redis.RDBClusterClient error=", err)
	}
	return rdsClusterClient
}

func RDBClusterClientInstance() *redis.ClusterClient {
	if rdsClusterClient != nil {
		return rdsClusterClient
	}
	rdsLock.Lock()
	defer rdsLock.Unlock()
	if rdsClusterClient != nil {
		return rdsClusterClient
	}
	return NewRDBClusterClient()
}
