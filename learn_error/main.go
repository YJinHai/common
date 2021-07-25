package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var Db *sql.DB
var dsnTest = GetEnv("DSN_TEST", "root:root@tcp(127.0.0.1:3306)/video")

// 基于该例子，个人认为最简单的方法是
// Sentinel Error + Wrap
// Sentinel Error 由db.go 公开，以避免造成导入循环
// "断言行为，而不是类型"这种处理方法，不适用于太简单的例子，导致增加代码量
// dao层的错误应该抛给上层，在dao层只记录相关错误信息，对于该错误的处理应该交给调用方，而且应该只有一次错误处理
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
