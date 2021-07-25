package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

var ErrNotFound = sql.ErrNoRows

func ConnectMySQL(dsn string) (*sql.DB, error) {
	return sql.Open("mysql", dsn)
}

// GetEnv get key environment variable if exist otherwise return defalutValue
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

type User struct {
	Id       string
	Username string
}

func QueryUserById(id string) (User, error) {
	var user User
	sqlS := "select id ,username from users where id = ?"
	row := Db.QueryRow(sqlS, id)
	err := row.Scan(&user.Id, &user.Username)
	if err != nil {
		return user, errors.Wrapf(sql.ErrNoRows, fmt.Sprintf("sql: %s error: %v", sqlS, err))
	}
	return user, nil
}
