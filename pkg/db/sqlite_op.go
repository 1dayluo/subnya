// sqlite_op.go
package db

import (
	"DomainMonitor/pkg/readconf"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db_1 = readconf.ReadSqlConfig("db_1")
var db_conn *sql.DB

func InitSqlClient() {
	//@title InitSqlClient
	//@param
	//Return
	var err error
	db_conn, err = sql.Open("sqlite3", db_1)
	if err != nil {
		panic(err)
	}
	// defer db_conn.Close()
	// 设计以下三类表单：
	// 当前监控 - 记录实时需要监控的域名
	// 本次新增 - 记录每次计划任务捕获到的新增域名
	// 本次删除 - 记录每次计划任务捕获到的减少域名
	domains_table_sql := `CREATE TABLE  IF NOT EXISTS domains(
		"DOMAIN"    TEXT     NOT NULL,
		"SUBDOMAIN"      TEXT  UNIQUE  NOT NULL,
		"UPDATETIME"     DATE     NOT NULL,
		"CHECKEDTIME" 	INT NOT NULL
	 );`
	addeddomain_table_sql := `CREATE TABLE IF NOT EXISTS  added_domains(
	"DOMAIN"    TEXT     NOT NULL,
	"SUBDOMAIN"      TEXT   UNIQUE  NOT NULL,
	"UPDATETIME"     DATE     NOT NULL,
	"CHECKEDTIME" 	INT NOT NULL
	);`
	deleteddomain_table_sql := `CREATE TABLE IF NOT EXISTS deleted_domains(
		"DOMAIN"    TEXT     NOT NULL,
		"SUBDOMAIN"      TEXT UNIQUE NOT NULL,
		"UPDATETIME"     DATE     NOT NULL,
		"CHECKEDTIME" 	INT NOT NULL
	);`
	create_tables_sql := []string{domains_table_sql, addeddomain_table_sql, deleteddomain_table_sql}
	for _, v := range create_tables_sql {
		query, err := db_conn.Prepare(v)
		if err != nil {
			fmt.Println("Err:", err)
			panic(err)
		}
		defer query.Exec()
		// defer db_conn.Close()
	}
}

func InsertAdded(domain string, subdomain string) {
	//@title InsertAdded
	//@param domain(string) subdomain(string) checked_time(int)
	//Return
	current_time := time.Now().Format("2006-01-02 15:04:05")
	tables := []string{"domains", "added_domains"}

	tx, err := db_conn.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	for _, table := range tables {
		stmt, err := tx.Prepare(fmt.Sprintf("INSERT OR REPLACE INTO %s (DOMAIN, SUBDOMAIN, UPDATETIME, CHECKEDTIME) VALUES (?, ?, ?, COALESCE((SELECT CHECKEDTIME FROM domains WHERE SUBDOMAIN = ?), 0) + 1)", table))
		if err != nil {
			panic(err)
		}
		defer stmt.Close()

		if _, err := stmt.Exec(domain, subdomain, current_time, subdomain); err != nil {
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
