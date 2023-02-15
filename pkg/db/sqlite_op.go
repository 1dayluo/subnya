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
	db_conn, err := sql.Open("sqlite3", db_1)
	if err != nil {
		panic(err)
	}
	defer db_conn.Close()
	// 设计以下三类表单：
	// 当前监控 - 记录实时需要监控的域名
	// 本次新增 - 记录每次计划任务捕获到的新增域名
	// 本次删除 - 记录每次计划任务捕获到的减少域名
	domains_table_sql := `CREATE TABLE domains(
		"DOMAIN"    TEXT     NOT NULL,
		"SUBDOMAIN"      TEXT    NOT NULL,
		"UPDATETIME"     DATE     NOT NULL,
		"CHECKEDTIME" 	INT NOT NULL
	 );`
	addeddomain_table_sql := `CREATE TABLE added_domains(
	"DOMAIN"    TEXT     NOT NULL,
	"SUBDOMAIN"      TEXT    NOT NULL,
	"UPDATETIME"     DATE     NOT NULL,
	"CHECKEDTIME" 	INT NOT NULL
	);`
	deleteddomain_table_sql := `CREATE TABLE added_domains(
		"DOMAIN"    TEXT     NOT NULL,
		"SUBDOMAIN"      TEXT    NOT NULL,
		"UPDATETIME"     DATE     NOT NULL,
		"CHECKEDTIME" 	INT NOT NULL
	);`
	create_tables_sql := []string{domains_table_sql, addeddomain_table_sql, deleteddomain_table_sql}
	for _, v := range create_tables_sql {
		query, err := db_conn.Prepare(v)
		if err != nil {
			panic(err)
		}
		query.Exec()
	}

	defer db_conn.Close()
}

func insert_added(domain string, subdomain string, checked_time int) {
	//@title insert_added
	//@param domain(string) subdomain(string) checked_time(int)
	//Return
	current_time := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db_conn.Prepare("INSERT INTO domains (DOMAIN, SUBDOMAIN, UPDATETIME, CHECKEDTIME) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(domain, subdomain, current_time, checked_time)
	if err != nil {
		panic(err)
	}
	fmt.Println("Row inserted successfully!")
}
