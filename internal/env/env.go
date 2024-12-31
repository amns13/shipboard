package env

import "github.com/redis/go-redis/v9"

type Env struct {
	Rdb *redis.Client
}

func GetEnv(redisUri string) (*Env, error) {
	opts, err := redis.ParseURL(redisUri)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	env := &Env{Rdb: client}
	return env, err
}
