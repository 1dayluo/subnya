package db

import (
	"DomainMonitor/pkg/readconf"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var db_1 = readconf.ReadSqlConfig("db_1")
var db_conn *sql.DB

func InitSqlClient() {
	db_conn, err := sql.Open("sqlite3", db_1)
	// 设计以下三类表单：
	// 当前监控 - 记录实时需要监控的域名
	// 本次新增 - 记录每次计划任务捕获到的新增域名
	// 本次删除 - 记录每次计划任务捕获到的减少域名
	if err != nil {
		panic(err)
	}
	defer db_conn.Close()
}

func insert_added() {
	//@title insert_added
	//@param
	//Return
	fmt.Println("插入数据到表：当前监控")
	fmt.Println("插入数据到表：本次新增")
}
