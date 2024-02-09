package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/phillipmugisa/go_user_api/data"
)

type Storage interface {
	GetUsers(params ...string) ([]*data.User, error)
	CreateUser(*data.User) ([]*data.User, error)
	DeleteUser(*data.User) error
	CompleteUserCheck(string) ([]*data.User, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}
	fmt.Println("Database Connection Successful...")
	return &PostgresStorage{
		db: db,
	}, nil
}

func initDB() (*sql.DB, error) {
	HOST := os.Getenv("POSTGRES_HOST")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB")
	PORT := os.Getenv("POSTGRES_PORT")
	username := os.Getenv("POSTGRES_USER")

	// make db connection
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, errors.New("Couldnot connect to database")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// create required tables in db: Users
func (s *PostgresStorage) SetUpDB() error {
	query := `CREATE TABLE IF NOT EXISTS Users (
		ID SERIAL PRIMARY KEY,
		UserName VARCHAR(255) NOT NULL UNIQUE,
		FirstName VARCHAR(255),
		LastName VARCHAR(255),
		Email VARCHAR(255) NOT NULL UNIQUE,
		Password VARCHAR(255) NOT NULL,
		Phone VARCHAR(20),
		Region VARCHAR(100) NOT NULL,
		UserGender VARCHAR(20),
		UserLanguage VARCHAR(50) NOT NULL,
		UserDateBirth DATE NOT NULL,
		Code VARCHAR(20) NOT NULL,
		check_status BOOLEAN DEFAULT FALSE
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) GetUsers(params ...string) ([]*data.User, error) {
	username := params[0]
	email := "email"
	if len(params) == 2 {
		email = params[1]
	}

	if username == "" && email == "" {
		return nil, errors.New("Pass either username or email")
	}

	query := `SELECT username, email, code FROM Users WHERE username = $1 OR email = $2`
	rows, err := s.db.Query(query, username, email)
	if err != nil {
		return nil, err
	}
	return scanUsers(rows)
}
func (s *PostgresStorage) CreateUser(u *data.User) ([]*data.User, error) {
	code := u.GenerateCode()

	query := `INSERT INTO Users (username, email, password, region, userLanguage, userDateBirth, firstName, lastName, phone, userGender, Code)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`

	_, err := s.db.Query(
		query,
		u.UserName,
		u.Email,
		u.Password,
		u.Region,
		u.UserLanguage,
		u.UserDateBirth.Format("2006-01-02"),
		u.FirstName,
		u.LastName,
		u.Phone,
		u.UserGender,
		code,
	)
	if err != nil {
		return nil, err
	}

	return s.GetUsers(u.UserName)
}
func (s *PostgresStorage) CompleteUserCheck(username string) ([]*data.User, error) {
	query := `UPDATE Users SET check_status = TRUE  WHERE username = $1`
	_, err := s.db.Exec(query, username)
	if err != nil {
		return nil, err
	}
	return s.GetUsers(username)
}

func (s *PostgresStorage) DeleteUser(u *data.User) error {
	query := `DELETE FROM Users WHERE username = $1`
	_, err := s.db.Exec(query, u.UserName)
	if err != nil {
		return err
	}
	return nil
}

func scanUsers(rows *sql.Rows) ([]*data.User, error) {
	users := []*data.User{}
	for rows.Next() {
		user := new(data.User)
		err := rows.Scan(
			&user.UserName,
			&user.Email,
			&user.Code,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
