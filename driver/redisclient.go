package driver

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	redis "gopkg.in/redis.v3"
)

//NewRedisConnector get 1 new redis client
func NewRedisConnector(rhost string) (*redis.Client, error) {
	//get handle
	var err error
	redisCache := redis.NewClient(&redis.Options{
		Addr:     rhost,
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 4000,
	})
	//wait till have a good conn
	for i := 0; i <= 100; i++ {
		_, err := redisCache.Ping().Result()
		if err != nil {
			log.Println("Redis:", err)
		} else {
			log.Println("Redis: Connected.")
			break
		}
		time.Sleep(time.Second * 3)
	}
	if err != nil {
		return nil, err
	}
	return redisCache, nil
}
