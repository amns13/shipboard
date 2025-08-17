package conf

import (
	"context"
	"html/template"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Env struct {
	Db        *pgxpool.Pool
	Rdb       *redis.Client
	Templates *template.Template
	Logger    *log.Logger
}

func LoadEnv(postgresUri string, redisUri string, templates []string) (*Env, error) {

	opts, err := redis.ParseURL(redisUri)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(opts)

	dbPool, err := pgxpool.New(context.Background(), postgresUri)
	if err != nil {
		return nil, err
	}

	tmpls := template.Must(template.ParseFiles(templates...))

	logger := log.Default()
	logger.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)

	env := &Env{Db: dbPool, Rdb: redisClient, Templates: tmpls, Logger: logger}
	return env, nil
}
