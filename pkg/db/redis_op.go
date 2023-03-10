// redis_op.go
package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/1dayluo/subnya/pkg/logutil"
	"github.com/1dayluo/subnya/pkg/readconf"

	"github.com/redis/go-redis/v9"
)

var db_link = readconf.ReadRedisConfig("addr")
var db_pass = readconf.ReadRedisConfig("password")
var db_db, _ = strconv.Atoi(readconf.ReadRedisConfig("db"))
var (
	rdb *redis.Client
)
var ctx = context.Background()

func init() {
	if err := logutil.Init(); err != nil {
		logutil.Logf("Failed to initialize logger: %v", err)
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     db_link,
		Password: db_pass,
		DB:       db_db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	rdb.Del(ctx, "fmd5test") // use For test
	if err != nil {
		logutil.Logf("[ERR]:", err)
	}
}

func SetMd5InDB(fileMd5 string) {
	//@title SetMd5InDB
	//@param
	//Return Nil
	err := rdb.SAdd(ctx, "fmd5", fileMd5, 0).Err()
	if err != nil {
		logutil.Logf("[Err]:error in redis:  %v", err)
	}
}

func UpdateFileMd5(fileName string, fileMd5 string) {
	//@title UpdateFileMd5
	//@param
	//Return Nil
	rdb.HSet(ctx, "file_md5", fileName, fileMd5)

}

func SearchFileMd5(fileName string) string {
	//@title SearchFileMd5
	//@param
	//Return Nil
	md5, err := rdb.HGet(ctx, "file_md5", fileName).Result()
	if err != nil {
		fmt.Println("\tError:", err)
	}

	return md5
}

func DelMd5InDB(md5 string) {
	//@title DelMd5InDB
	//@param
	//Return
	_, err := rdb.SRem(ctx, "fmd5", md5).Result()
	if err != nil {
		fmt.Println("\tError:", err)
	}
}

func CheckMd5InDB(fileMd5 string) bool {
	//@title checkFileMd5
	//@param
	//Return Nil

	exists, err := rdb.SIsMember(ctx, "fmd5", fileMd5).Result()
	if err != nil {
		return false
	}
	if exists {
		// fmt.Println("EXISTS!")
		return true
	}
	// fmt.Println("NOT EXISTS")
	return false

}

func InsertNewFindMd5(fname string, fmd5 string) {
	//@title InsertNewFindMd5:
	//@param
	//Return
	SetMd5InDB(fmd5)
	UpdateFileMd5(fname, fmd5)

}
