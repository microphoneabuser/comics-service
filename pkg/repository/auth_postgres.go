package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	serv "github.com/microphoneabuser/comics_service"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user serv.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, login, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	
	row := r.db.QueryRow(query, user.Name, user.Login, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(login, password string) (serv.User, error) {
	var user serv.User
	query := fmt.Sprintf("SELECT id, name, login FROM %s WHERE login=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, login, password)

	return user, err
}

func (r *AuthPostgres) GetUserById(id int) (serv.User, error) {
	var user serv.User
	query := fmt.Sprintf("SELECT id, name, login FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)

	return user, err
}
