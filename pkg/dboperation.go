package pkg

import (
	"DomainMonitor/pkg/util"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var db_link = util.ReadRedisConfig("addr")
var db_pass = util.ReadRedisConfig("password")
var db_db, _ = strconv.Atoi(util.ReadRedisConfig("db"))
var rdb = redis.NewClient(&redis.Options{
	Addr:     db_link,
	Password: db_pass,
	DB:       db_db,
})

func checkMd5InDB(fileMd5 string) {
	//@title checkFileMd5
	//@param
	//Return Nil

}
