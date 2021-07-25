package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var Db *sql.DB
var dsnTest = GetEnv("DSN_TEST", "root:root@tcp(127.0.0.1:3306)/video")

//
func main() {
	var err error
	Db, err = ConnectMySQL(dsnTest)
	if err != nil {
		panic(err)
	}

	user, err2 := QueryUserById("1")
	if err2 != nil && !errors.Is(err2, ErrNotFound) {
		// log 日志记录打印的内容
		fmt.Println(fmt.Sprintf("%+v", err2))

		// return 给客户端的内容
		fmt.Println("查询出错")
		return
	}
	fmt.Println(user)
}
