package db

import "github.com/go-redis/redis"

var redisDb *redis.Client

// 初始化连接
func initClient() (err error) {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     "47.111.175.222:6379",
		Password: "fuckyou!",
		DB:       0,
	})

	_, err = redisDb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func GetRedisClient() (*redis.Client, error) {
	err := initClient()
	if err != nil {
		return nil, err
	}
	return redisDb, nil

}
