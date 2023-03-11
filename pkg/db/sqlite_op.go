// sqlite_op.go
package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/1dayluo/subnya/pkg/logutil"
	"github.com/1dayluo/subnya/pkg/readconf"

	_ "github.com/mattn/go-sqlite3"
)

var db_1 = readconf.ReadSqlConfig("db_1")
var db_conn *sql.DB

type SubdomainInfos struct {
	Domain      string
	Subdomain   string
	Updatetime  string
	Checkedtime int
	IFON        int
	STATUS      int
}

func setDefault() {
	if _, err := os.Stat(db_1); os.IsNotExist(err) {
		homedir, _ := os.UserHomeDir()
		dbPath := fmt.Sprintf("%v/.config/subnya/db/", homedir)
		db_1 = readconf.SetSqliteConfig("db_1", dbPath+"monitor.db")
		os.MkdirAll(dbPath, os.ModePerm)
	}
}
func init() {
	//@title InitSqlClient
	//@param
	//Return
	if err := logutil.Init(); err != nil {
		logutil.Logf("Failed to initialize logger: %v", err)
	}
	setDefault()

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
		"CHECKEDTIME" 	INT NOT NULL,
		"IFON"		INT NOT NULL,
		"STATUS"		INT
	 );`
	addeddomain_table_sql := `CREATE  TABLE IF NOT EXISTS  added_domains(
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
			logutil.Logf("[Err]:error in sql:", err)
			panic(err)
		}
		defer query.Exec()
		// defer db_conn.Close()
	}
}

func Test(domain string, subdomain string) {
	current_time := time.Now().Format("2006-01-02 15:04:05")
	tx, err := db_conn.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		err = tx.Commit()
	}()
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO domains (DOMAIN, SUBDOMAIN, UPDATETIME, CHECKEDTIME, IFON) VALUES (?, ?, ?, COALESCE((SELECT CHECKEDTIME FROM domains WHERE SUBDOMAIN = ?), 0) + 1, 1)")
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(domain, subdomain, current_time, subdomain); err != nil {
		return
	}

}
func AddMonitor(domain string, subdomain string, status int) (err error) {
	//@title AddMonitor
	//@param domain(string) subdomain(string) checked_time(int)
	//Return
	current_time := time.Now().Format("2006-01-02 15:04:05")
	tx, err := db_conn.Begin()
	if err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		err = tx.Commit()
	}()
	sql_operation := []string{
		fmt.Sprintf("INSERT OR REPLACE INTO %s (DOMAIN, SUBDOMAIN, UPDATETIME, CHECKEDTIME, IFON, STATUS) VALUES (?, ?, ?, COALESCE((SELECT CHECKEDTIME FROM domains WHERE SUBDOMAIN = ?), 0) + 1, 1, ?)", "domains"),
		fmt.Sprintf("INSERT OR REPLACE INTO %s (DOMAIN, SUBDOMAIN, UPDATETIME, CHECKEDTIME) VALUES (?, ?, ?, COALESCE((SELECT CHECKEDTIME FROM domains WHERE SUBDOMAIN = ?), 0 ))", "added_domains"),
	}
	stmt_1, err := tx.Prepare(sql_operation[0])
	if err != nil {
		logutil.Logf("[ERROR]error when update tables:%v", err)
		panic(err)
	}
	defer stmt_1.Close()

	if _, err := stmt_1.Exec(domain, subdomain, current_time, subdomain, status); err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
		// return
	}

	stmt_2, err := tx.Prepare(sql_operation[1])
	if err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
	}
	defer stmt_2.Close()

	if _, err := stmt_2.Exec(domain, subdomain, current_time, subdomain); err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
		// return
	}
	return
}

func DeleteMonitor(domain string, subdomain string, status int) (err error) {
	//@title DeleteMonitor
	//@param domain(string) subdomain(string) checked_time(int)
	//Return
	current_time := time.Now().Format("2006-01-02 15:04:05")

	tx, err := db_conn.Begin()
	if err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		err = tx.Commit()
	}()
	sql_operation := []string{
		fmt.Sprintf("INSERT OR REPLACE INTO %s (DOMAIN, SUBDOMAIN, UPDATETIME, CHECKEDTIME, IFON, STATUS) VALUES (?, ?, ?, COALESCE((SELECT CHECKEDTIME FROM domains WHERE SUBDOMAIN = ?), 0) + 1, 0, ?)", "domains"),
		fmt.Sprintf("INSERT OR REPLACE INTO %s (DOMAIN, SUBDOMAIN, UPDATETIME, CHECKEDTIME) VALUES (?, ?, ?, COALESCE((SELECT CHECKEDTIME FROM domains WHERE SUBDOMAIN = ?), 0) )", "deleted_domains"),
	}
	stmt_1, err := tx.Prepare(sql_operation[0])
	if err != nil {
		logutil.Logf("[ERROR]error when update tables:%v", err)
		panic(err)
	}
	defer stmt_1.Close()
	if _, err := stmt_1.Exec(domain, subdomain, current_time, subdomain, status); err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
		// return
	}

	stmt_2, err := tx.Prepare(sql_operation[1])
	if err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
	}
	defer stmt_2.Close()
	if _, err := stmt_2.Exec(domain, subdomain, current_time, subdomain); err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
		// return
	}
	return
}

func GetSubDomianInfo(subdomain string) (dinfos []SubdomainInfos) {
	//@title GetSubDomianInfo
	//@param
	//Return
	sql_operation := "SELECT * FROM domains WHERE subdomain = ?"
	stmt, err := db_conn.Prepare(sql_operation)
	if err != nil {
		logutil.Logf("[ERROR]error when get infos from tables:%v", err)
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(subdomain)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var domain, subdomain, updatetime string
		var checked_time, ifon, status int
		if err := rows.Scan(&domain, &subdomain, &updatetime, &checked_time, &ifon, &status); err != nil {
			panic(err)
		}
		info := SubdomainInfos{
			Domain:      domain,
			Subdomain:   subdomain,
			Checkedtime: checked_time,
			Updatetime:  updatetime,
			IFON:        ifon,
			STATUS:      status,
		}
		dinfos = append(dinfos, info)
		// fmt.Printf("Domain: %s, Subdomain: %s, Updatetime: %s, Checkedtime: %d\n", _domain, _subdomain, _updatetime, _checked_time)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return
}

func Getdomains() (domains []string) {
	//@title getDomains
	//@param
	//Return
	sql_op := "SELECT DOMAIN FROM domains GROUP BY DOMAIN"
	stmt, err := db_conn.Prepare(sql_op)
	if err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			logutil.Logf("[Err]:error in sql:  %v", err)
			panic(err)
		}
		domains = append(domains, domain)
		// fmt.Printf("Domain: %s, Subdomain: %s, Updatetime: %s, Checkedtime: %d\n", _domain, _subdomain, _updatetime, _checked_time)
	}
	if err := rows.Err(); err != nil {
		logutil.Logf("[ERROR]%v", err)
		panic(err)
	}
	return
}

func GetMonitoredSub(domain string) (domains []string) {
	//@title getDomains
	//@param
	//Return
	sql_op := "SELECT SUBDOMAIN FROM domains  WHERE IFON = 1 AND DOMAIN = ?"
	stmt, err := db_conn.Prepare(sql_op)
	if err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(domain)
	if err != nil {
		logutil.Logf("[Err]:error in sql:  %v", err)
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			panic(err)
		}
		domains = append(domains, domain)
		// fmt.Printf("Domain: %s, Subdomain: %s, Updatetime: %s, Checkedtime: %d\n", _domain, _subdomain, _updatetime, _checked_time)
	}
	if err := rows.Err(); err != nil {
		logutil.Logf("[ERROR]%v", err)
		panic(err)
	}
	return
}
