package repo

import (
	"context"
	"database/sql"

	"github.com/bipen2001/go-user-assignment-go/internal/entity"
	"github.com/bipen2001/go-user-assignment-go/internal/service/user/model"
)

type repository struct {
	db *sql.DB
}

func NewRepository(dialect, dsn string, idleConn, maxConn int) (model.Repository, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	return &repository{db}, nil
}

func (r *repository) GetById(ctx context.Context, id string) (*entity.User, error) {
	var user *entity.User
	row := r.db.QueryRow(`SELECT * FROM users WHERE value=$1`, id)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Dob)

	if err != nil {
		if err == sql.ErrNoRows {

			return nil, err
		}
		return nil, err
	}

	r.db.Close()

	return user, nil
}

func (r *repository) Get(ctx context.Context) ([]entity.User, error) {
	var user entity.User

	rows, err := r.db.Query(`SELECT id,firstName,lastName,email,dob FROM users`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []entity.User
	for rows.Next() {
		rows.Scan(&user)
		users = append(users, user)
	}
	r.db.Close()
	return users, nil
}

func (r *repository) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	var u entity.User

	statement := `INSERT INTO users (firstName,lastName,email,password,dob) VALUES($1,$2,$3,$4,$5) RETURNING *`
	rows, err := r.db.Query(statement, user.FirstName, user.LastName, user.Email, user.Password, user.Dob)
	rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.Dob)
	if err != nil {
		return nil, err
	}
	r.db.Close()

	return &u, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {

	statement := `DELETE FROM users where id=$1`
	_, err := r.db.Exec(statement, id)

	if err != nil {
		return err
	}
	r.db.Close()

	return nil
}
func (r *repository) Update(ctx context.Context, user entity.User) (*entity.User, error) {
	var u entity.User

	statement := `UPDATE users SET first_name = $2, last_name = $3 WHERE id = $1`
	rows, err := r.db.Query(statement, user.FirstName, user.LastName, user.Id)
	rows.Scan(&u.FirstName, &u.LastName, &u.Email, &u.Dob)
	if err != nil {
		return nil, err
	}
	r.db.Close()

	return &u, nil
}
