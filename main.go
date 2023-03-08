/*
 * @Author: 1dayluo
 * @Date: 2023-02-08 09:55:51
 * @LastEditTime: 2023-03-08 09:03:10
 */
package main

import (
	"DomainMonitor/pkg/cmd"
	// redis "DomainMonitor/pkg/db"
	sqlite "DomainMonitor/pkg/db"
	"DomainMonitor/pkg/io"
	"DomainMonitor/pkg/logutil"
	"DomainMonitor/pkg/output"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/alexflint/go-arg"
)

type args struct {
	// -u 查找文件md5的更新，有更新则会单独跑一次数据
	// -r 对数据库内的监控文件进行内容读取(不会对文件的更新进行追踪），并查找每个域名下可能的子域名。（最后存储到数据库中/验活）
	// -output 输出本次更新统计结果的文件|默认输出在终端下
	UPDATE bool `arg:"-u,--update" help:"Check update in monitor"`
	RUN    bool `arg:"-r,--run" help:"start subdomain finder and update data(include response status code) in sqlite"`
	OUT    bool `arg:"--output" helo:"output to file defined in default settings`
}

type CacheDomain struct {
	Domain    string
	Subdomain string
}

func formaturl(url string) (furl string) {
	/**
	 * @description: format url
	 * @return {*}
	 */
	domain := strings.Split(url, "//")
	if len(domain) == 1 {
		furl = "http://" + strings.Join(domain, "")
	} else {
		furl = url
	}
	return
}
func aliveCheck(url string) (bool, int) {
	/**
	 * @description: send requests to url,and check response code
	 * @return {*}
	 */
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

func difference(a, b []string) []string {
	/**
	 * @description: Compare the difference between a and b, and return as a slice
	 * @return {*}
	 */
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
	/**
	 * @description: Find out the same element between a and b, and return it as a slice
	 * @return {*}
	 */
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
	/**
	 * @description: Compare return_domains with the domains in the database to find out the new and reduced parts relative to the data in the database
	 * @return {*}
	 */
	same_domain := intersect(return_domains, db_domains)
	if len(return_domains) > len(db_domains) {
		added_domains = difference(return_domains, same_domain)
	}
	if len(return_domains) < len(db_domains) {
		deled_domains = difference(db_domains, same_domain)
	}
	return added_domains, deled_domains

}
func upgradeAddSubdomainSQL(domain string, subdomains []string) {
	/**
	 * @description: Update the newly added subdomain  detected under the domain  to the database
	 * @return {*}
	 */
	for _, subdomain := range subdomains {
		// res, code := aliveCheck(subdomain)
		// fmt.Println(res, code, subdomain)
		sqlite.AddMonitor(domain, subdomain, -1)
	}
}

func upgradeDelSubdomainSQL(domain string, subdomains []string) {
	/**
	 * @description: Update the detected reduced subdomain  under the domain  to the database
	 * @return {*}
	 */
	for _, subdomain := range subdomains {
		// res, code := aliveCheck(subdomain)
		// fmt.Println(res, code, subdomain)
		sqlite.DeleteMonitor(domain, subdomain, -1)
	}
}
func Update(domains []string) (results map[string]output.ResultOutput) {
	/**
	 * @description: Update the subdomain under the domain
	 * @return {*}
	 */
	var wg sync.WaitGroup
	var mu sync.Mutex
	results = make(map[string]output.ResultOutput)
	domainCH := make(chan string)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range domainCH {
				mu.Lock()
				monitoredSub := sqlite.GetMonitoredSub(domain)
				mu.Unlock()
				subdomains := strings.Split(cmd.FindStart(domain), "\n")
				added, deled := get_changed(subdomains, monitoredSub)
				if added != nil {
					mu.Lock()
					upgradeAddSubdomainSQL(domain, subdomains)
					results[domain] = output.ResultOutput{All: subdomains, Added: subdomains, Deled: nil}
					// resuts = append(resuts, ResultOutput{Domain: domain, Added: subdomains, Deled: nil})
					mu.Unlock()
				}
				if deled != nil {
					mu.Lock()
					upgradeDelSubdomainSQL(domain, subdomains)
					results[domain] = output.ResultOutput{All: subdomains, Added: nil, Deled: subdomains}
					// resuts = append(resuts, ResultOutput{Domain: domain, Added: nil, Deled: subdomains})
					mu.Unlock()
				}
			}
		}()
	}

	for _, domain := range domains {
		domainCH <- domain
	}

	close(domainCH)
	wg.Wait()
	return

}

func UpdateAndCheck(domains []string) (results map[string]output.ResultOutput) {
	/**
	 * @description: Update the subdomain under the domain and check its response code
	 * @return {*}
	 */

	var wg sync.WaitGroup
	var mu sync.Mutex
	results = make(map[string]output.ResultOutput)
	subdomainsCH := make(chan CacheDomain)
	resultsCH := make(chan output.ResultOutput)

	// Spawn worker goroutines
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for subdomainInfo := range subdomainsCH {
				_, code := aliveCheck(subdomainInfo.Subdomain)
				mu.Lock()
				results[subdomainInfo.Domain].Code[subdomainInfo.Subdomain] = code
				sqlite.AddMonitor(subdomainInfo.Domain, subdomainInfo.Subdomain, code)
				mu.Unlock()
			}
		}()
	}

	// Spawn separate goroutine to send domains to the channel
	go func() {
		for _, domain := range domains {
			subdomains := strings.Split(cmd.FindStart(domain), "\n")
			monitoredSub := sqlite.GetMonitoredSub(domain)
			added, deled := get_changed(subdomains, monitoredSub)
			mu.Lock()
			results[domain] = output.ResultOutput{All: subdomains, Added: added, Deled: deled, Code: map[string]int{}}
			mu.Unlock()
			for _, subdomain := range subdomains {
				subdomainsCH <- CacheDomain{Domain: domain, Subdomain: subdomain}
			}
		}
		close(subdomainsCH)
	}()

	wg.Wait()
	close(resultsCH)
	return
}

func main() {

	if err := logutil.Init(); err != nil {
		logutil.Logf("Failed to initialize logger: %v", err)
	}
	var results map[string]output.ResultOutput

	var args args
	arg.MustParse(&args)

	// 检查监控对象，查看是否有新增监控对象，并对新增监控对象进行子域名查询/更新
	if args.UPDATE {
		files := io.SearchAndUpdateMd5()
		fmt.Printf("[Info]New find in files: %v", files)
		for _, file := range files {
			lines := io.ReadFileContent(file)
			// TODO: 检查取消监控的域名
			results = Update(lines)
			output.OutResult(results, os.Stdout, true)
		}
	}
	// 运行sqlite内记录的域名，查看是否有子域名更新
	if args.RUN {
		domains := sqlite.Getdomains()
		results = UpdateAndCheck(domains)
		output.OutResult(results, os.Stdout, true)

	}
	if args.OUT {
		fmt.Println(output.OutFile)
		file, err := os.Create(output.OutFile)
		if err != nil {
			logutil.Logf("Output failed: %v", err)
		}
		output.OutResult(results, file, false)
	}

	// output.OutResult(test)

}
