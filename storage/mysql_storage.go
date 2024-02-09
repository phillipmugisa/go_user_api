package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/phillipmugisa/go_user_api/data"
)

type Storage interface {
	GetUser(string) ([]*data.User, error)
	CreateUser(*data.User) ([]*data.User, error)
	CompleteUserCheck(string) ([]*data.User, error)
}

type MySqlStorage struct {
	db *sql.DB
}

func NewMySqlStorage() (*MySqlStorage, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}
	return &MySqlStorage{
		db: db,
	}, nil
}

func initDB() (*sql.DB, error) {
	// make db connection
	dbUrl := GetConnectionString()

	fmt.Println("dbUrl: ", dbUrl)
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, errors.New("Couldnot connect to database")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *MySqlStorage) SetUpDB() error {
	query := `CREATE TABLE IF NOT EXISTS Users (
		ID INT AUTO_INCREMENT PRIMARY KEY,
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

func (s *MySqlStorage) GetUser(username string) ([]*data.User, error) {
	query := `SELECT username, email, code FROM Users WHERE username = ?`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	return scanUsers(rows)
}
func (s *MySqlStorage) CreateUser(u *data.User) ([]*data.User, error) {
	code := u.GenerateCode()

	query := `INSERT INTO Users (username, email, password, region, userLanguage, userDateBirth, firstName, lastName, phone, userGender, Code)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

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

	return s.GetUser(u.UserName)
}
func (s *MySqlStorage) CompleteUserCheck(username string) ([]*data.User, error) {
	query := `UPDATE Users SET check_status = TRUE  WHERE username = ?`
	_, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	return s.GetUser(username)
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

func GetConnectionString() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("DB_PASS")
	if password == "" {
		password = "@root"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "apiserver"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)
}
