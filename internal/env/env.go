package env

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Env struct {
	Db  *pgxpool.Pool
	Rdb *redis.Client
}

func LoadEnv(postgresUri string, redisUri string) (*Env, error) {

	opts, err := redis.ParseURL(redisUri)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(opts)

	dbPool, err := pgxpool.New(context.Background(), postgresUri)
	if err != nil {
		return nil, err
	}

	env := &Env{Db: dbPool, Rdb: redisClient}
	return env, nil
}
