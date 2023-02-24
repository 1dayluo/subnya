package main

import (
	"DomainMonitor/pkg/cmd"
	redis "DomainMonitor/pkg/db"
	sqlite "DomainMonitor/pkg/db"
	"DomainMonitor/pkg/io"
	"DomainMonitor/pkg/readconf"
	"fmt"
	"net/http"
	"strings"
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

func formaturl(url string) (furl string) {
	domain := strings.Split(url, "//")
	if len(domain) == 1 {
		furl = "http://" + strings.Join(domain, "")
	} else {
		furl = url
	}
	return
}
func aliveCheck(url string) (bool, int) {
	//@title aliveCheck
	//@param
	//Return bool
	// timeout := time.Duration(2*time.Second)
	url = formaturl(url)
	resp, err := http.Get(url)
	if err != nil {
		// fmt.Println(err)
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
func difference(a, b []string) []string {
	//@title difference
	//@param
	//Return
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
func intersect(slice1, slice2 []string) []string {
	var intersect []string
	for _, element1 := range slice1 {
		for _, element2 := range slice2 {
			if element1 == element2 {
				intersect = append(intersect, element1)
			}
		}
	}
	return intersect //return slice after intersection
}
func get_changed(return_domains, db_domains []string) (added_domains []string, deled_domains []string) {
	//@title get_changed
	//@param
	//Return
	same_domain := intersect(return_domains, db_domains)
	if len(return_domains) > len(db_domains) {
		added_domains = difference(return_domains, same_domain)
	} else {
		deled_domains = difference(db_domains, same_domain)
	}
	return added_domains, deled_domains

}
func upgradeSubdomainSQL(domain string, subdomains []string) {
	//@title insertFinder
	//@param
	//Return
	for _, subdomain := range subdomains {
		// res, code := aliveCheck(subdomain)
		// fmt.Println(res, code, subdomain)
		sqlite.AddMonitor(domain, subdomain, -1)
	}
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
				subdomains := strings.Split(cmd.FindStart(domain), "\n")
				upgradeSubdomainSQL(domain, subdomains)

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
	if err := redis.InitClient(); err != nil {
		panic(err)
	}
	var args args
	sqlite.InitSqlClient()
	arg.MustParse(&args)

	if args.UPDATE {
		files := searchAndUpdateMd5()
		fmt.Printf("\t[Info]New find in files: %v", files)
		for _, file := range files {
			fmt.Println(file)
			lines := io.ReadFileContent(file)
			scanSubdomain(lines)
		}

	}
	if args.RUN {
		domains := sqlite.Getdomains()
		scanSubdomain(domains)
	}
	if args.OUTPUT != nil {
		fmt.Println("some code")
	}

	// fmt.Printf("\t[Info]New find in files: %v", files)
	//
	// sqlite.Test("abc.com", "xyz.abc.com")
	// sqlite.AddMonitor("abc.com", "xyz.abc.com")

}
