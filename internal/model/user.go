package model

import (
	"context"
	"net/mail"
	"time"

	"github.com/amns13/shipboard/internal/conf"
	"github.com/jackc/pgx/v5"
)

type UserCreator struct {
	Name         string `db:"name"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

type User struct {
	Id        int32     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UserCreator
}

// For now, we return everything and use returning. If for some hypothetical reason,
// we need to scale, that can be done easily.
const insertUserQuery = `
INSERT INTO users (name, email, password_hash)
VALUES (@name, @email, @password_hash)
RETURNING *;
`

const userExistsQuery = `
SELECT EXISTS(SELECT 1 FROM users WHERE email = @email);
`

// TODO: Some fields are unneeded. Keeping them for now.
const userSelectFromEmailQuery = `
SELECT id, email, password_hash, name, created_at
FROM users
WHERE email = @email;
`

const userSelectFromIdQuery = `
SELECT id, email, password_hash, name, created_at
FROM users
WHERE id = @id;
`

func (usr *UserCreator) Create(env *conf.Env) (*User, error) {

	args := pgx.NamedArgs{
		"name":          usr.Name,
		"email":         usr.Email,
		"password_hash": usr.PasswordHash,
	}
	// Returned error will be handled while parsing returnedRows
	returnedRows, _ := env.Db.Query(context.Background(), insertUserQuery, args)
	user, err := pgx.CollectOneRow(returnedRows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UserExists(env *conf.Env, email *mail.Address) (bool, error) {
	args := pgx.NamedArgs{
		"email": email,
	}
	var exists bool
	err := env.Db.QueryRow(context.Background(), userExistsQuery, args).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetUserByEmail(env *conf.Env, email string) (*User, error) {
	args := pgx.NamedArgs{
		"email": email,
	}
	var user *User
	// Returned error will be handled while parsing returnedRows
	returnedRows, _ := env.Db.Query(context.Background(), userSelectFromEmailQuery, args)
	user, err := pgx.CollectOneRow(returnedRows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(env *conf.Env, id int32) (*User, error) {
	args := pgx.NamedArgs{
		"id": id,
	}
	var user *User
	// Returned error will be handled while parsing returnedRows
	returnedRows, _ := env.Db.Query(context.Background(), userSelectFromIdQuery, args)
	user, err := pgx.CollectOneRow(returnedRows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, err
	}
	return user, nil
}
