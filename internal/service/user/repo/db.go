package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bipen2001/go-user-assignment-go/internal/entity"
	"github.com/bipen2001/go-user-assignment-go/internal/service/user/model"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type repository struct {
	db *sql.DB
}

func NewRepository(dialect, dsn string) (model.Repository, error) {

	db, err := sql.Open(dialect, dsn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return &repository{db}, nil
}

func (r *repository) GetById(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User

	row := r.db.QueryRow(`SELECT id,firstName,lastName,email,dob FROM users WHERE id=$1`, id)

	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Dob)

	if err != nil {

		if err == sql.ErrNoRows {

			return nil, errors.New("No user With this id exists" + err.Error())
		}

		return nil, err
	}

	return &user, nil
}

func (r *repository) Get(ctx context.Context, queryParams entity.QueryParams, pass bool) ([]entity.Response, error) {

	var user entity.Response

	statement := createSearchQuery(
		queryParams.Archived,
		queryParams.FirstName,
		queryParams.LastName,
		queryParams.Email,
		queryParams.Sort,
		queryParams.Order,
		pass,
	)

	rows, err := r.db.Query(statement)

	if err != nil {

		return nil, errors.New("No records found" + err.Error())
	}

	defer rows.Close()

	var users []entity.Response

	for rows.Next() {

		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Dob)

		if err != nil {

			return nil, errors.New("No Users found" + err.Error())
		}
		user.Age = calculateAge(user.Dob)

		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, errors.New("no records found")

	}

	return users, nil
}

func (r *repository) Create(ctx context.Context, user entity.User) (*entity.Response, error) {
	var u entity.Response

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("password encyption failed " + err.Error())

	}

	statement := `INSERT INTO users (
		id, firstName, lastName, email, password, 
		dob
	  ) 
	  VALUES 
		($1, $2, $3, $4, $5, $6) RETURNING id, 
		firstName, 
		lastName, 
		email, 
		dob
	  `

	row := r.db.QueryRow(statement,
		uuid.NewV4(),
		user.FirstName,
		user.LastName,
		user.Email,
		string(hashedPassword),
		user.Dob,
	)

	err = row.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.Dob)

	if err != nil {

		return nil, errors.New("validation error " + err.Error())
	}
	u.Age = calculateAge(u.Dob)

	return &u, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {

	statement := `UPDATE users SET archived='TRUE' WHERE id=$1`

	_, err := r.db.Exec(statement, id)

	if err != nil {

		return errors.New("Could not delete User " + err.Error())
	}

	return nil
}

func (r *repository) Update(ctx context.Context, id string, user entity.UpdateUser) (*entity.User, error) {

	var u entity.User

	statement := `UPDATE users SET updated_at='NOW()', ` + user.Update() + ` WHERE id = $1 Returning id,firstName,lastName,email,dob`

	row := r.db.QueryRow(statement, id)

	err := row.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.Dob)
	if err != nil {

		return nil, errors.New("could not update User:" + err.Error())
	}

	return &u, nil
}

//db functions

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)

	if err != nil {
		fmt.Println(err)
		return err
	}

	statement := `
	CREATE TABLE IF NOT EXISTS users(
		id UUID PRIMARY KEY ,
		firstName VARCHAR(30) NOT NULL,
		lastName VARCHAR(30) NOT NULL,
		email VARCHAR NOT NULL UNIQUE,
		password VARCHAR(20) NOT NULL,
		dob timestamp NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		archived BOOLEAN NOT NULL DEFAULT FALSE  
	)
`
	_, err = db.Exec(statement)

	if err != nil {
		return err
	}

	return db.Close()
}

func DropDb(driverName, dataSource, dbname string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DROP DATABASE IF EXISTS ` + dbname)
	if err != nil {
		return err
	}
	return nil
}

func calculateAge(t *time.Time) int {

	return int(time.Since(*t).Hours()) / (365 * 24)

}

func createSearchQuery(
	archived string,
	firstName string,
	lastName string,
	email string,
	// created_at *time.Time,
	// last_accessed_at *time.Time,
	// updated_at *time.Time,
	sort string,
	order string,
	pass bool,

) string {

	var arr []string
	var where string
	var orderBy string
	var password string
	if pass {
		password = ", password "
	}
	if strings.ToUpper(archived) != "TRUE" {
		arr = append(arr, "archived=FALSE")
	}

	if firstName != "" {
		arr = append(arr, " firstName LIKE '%"+firstName+"%'")

	}
	if lastName != "" {
		arr = append(arr, " lastname LIKE '%"+lastName+"%'")

	}
	if email != "" {
		arr = append(arr, " email LIKE '%"+email+"%'")

	}
	if len(arr) > 0 {

		where = ` WHERE`

	} else {
		where = ""
	}
	if order == "" {
		order = "ASC"
	}
	if sort != "" {
		orderBy = "ORDER BY " + sort + " " + strings.ToUpper(order)
	}

	query := `SELECT id,firstName,lastName,email,dob ` + password + ` FROM users` + where + " " + strings.Join(arr, " AND ") + " " + orderBy
	fmt.Println(query)
	return query
}
