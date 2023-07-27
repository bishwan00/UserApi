// initializers/redis.go

package initializers

import (
    "context"
    "github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func SetupRedis() {
    RDB = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Redis server address
        Password: "",               // No password set
        DB:       0,                // Use the default DB
    })

    _, err := RDB.Ping(Ctx).Result()
    if err != nil {
        panic(err)
    }
}
