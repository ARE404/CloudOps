package di

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Redis 连接失败: addr=%s err=%v", viper.GetString("redis.addr"), err))
	}
	return client
}
