package main

import (
	"DomainMonitor/pkg/cmd"
	redis "DomainMonitor/pkg/db"
	"DomainMonitor/pkg/io"
	"DomainMonitor/pkg/readconf"
	"fmt"
	"net/http"
	"sync"

	"github.com/alexflint/go-arg"
)

type args struct {
	// -u 查找文件md5的更新，有更新则会单独跑一次数据
	// -r 对数据库内的监控文件进行内容读取，并查找每个域名下可能的子域名。（最后存储到数据库中）
	// -output 输出本次更新统计结果的文件|默认输出在终端下
	UPDATE bool     `arg:"-u,--update" help:"Check update in monitor"`
	RUN    bool     `arg:"-r,--run" help:"start subdomain finder and update data in sqlite"`
	OUTPUT []string `arg:"--output"`
}

func InsertNewFindMd5(fname string, fmd5 string) {
	//@title InsertNewFindMd5:
	//@param
	//Return
	redis.SetMd5InDB(fmd5)
	redis.UpdateFileMd5(fname, fmd5)

}

func aliveCheck(url string) (bool, int) {
	//@title aliveCheck
	//@param
	//Return bool
	// timeout := time.Duration(2*time.Second)
	resp, err := http.Get(url)
	if err != nil {
		return false, -1
	} else {
		return true, resp.StatusCode
	}
}

func searchAndUpdateMd5() (newMonitorFiles []string) {
	//@title searchAndUpdateMd5
	//@param
	//Return newMonitorFIles []string (Files changed during this check)
	dirs := readconf.ReadMonitorDir()
	var dirInfo []map[string]string
	for _, dir := range dirs {
		dirInfo = append(dirInfo, io.ReadFromDir(dir))
		for _, finfos := range dirInfo {
			for fname, fmd5 := range finfos {
				fname = dir + "/" + fname
				if redis.CheckMd5InDB(fmd5) {
					// fmt.Println("\t[Info]Exists:", fname)
					history_md5 := redis.SearchFileMd5(fname)
					if history_md5 != fmd5 {
						InsertNewFindMd5(fname, fmd5)
						newMonitorFiles = append(newMonitorFiles, fname)
					}
				} else { //If fmd5 not in
					// fmt.Println("\t[Info]Not Exists:", fname)
					redis.UpdateFileMd5(fname, fmd5)
					InsertNewFindMd5(fname, fmd5)
					newMonitorFiles = append(newMonitorFiles, fname)
				}
			}
		}
	}

	return
}

func scanSubdomain(domains []string) {
	//@title scanSubdomain
	//@param
	//Return
	var wg sync.WaitGroup
	domainCH := make(chan string)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range domainCH {

				cmd.FindStart(domain)

			}
		}()
	}
	for _, domain := range domains {
		domainCH <- domain
	}

	close(domainCH)
	wg.Wait()
}

func main() {
	fmt.Println("Hello World!")
	redis.InitClient()
	var args args

	arg.MustParse(&args)

	if args.UPDATE {
		files := searchAndUpdateMd5()
		fmt.Printf("\t[Info]New find in files: %v", files)
		for _, file := range files {
			fmt.Println(file)
			// lines := io.ReadFileContent(file)

		}

	}
	if args.RUN {
		fmt.Println("some code")
	}
	if args.OUTPUT != nil {
		fmt.Println("some code")
	}

	// fmt.Printf("\t[Info]New find in files: %v", files)
	// sqlite.InitSqlClient()
	// sqlite.Test("abc.com", "xyz.abc.com")
	// sqlite.AddMonitor("abc.com", "xyz.abc.com")

}
