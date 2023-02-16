package main

import (
	redis "DomainMonitor/pkg/db"
	sqlite "DomainMonitor/pkg/db"
	"DomainMonitor/pkg/io"
	"DomainMonitor/pkg/readconf"
	"fmt"
	"net/http"
)

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
	//Return Nil
	dirs := readconf.ReadMonitorDir()
	var dirInfo []map[string]string
	for _, dir := range dirs {
		dirInfo = append(dirInfo, io.ReadFromDir(dir))
	}
	for _, finfos := range dirInfo {
		for fname, fmd5 := range finfos {
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
	return
}
func main() {
	fmt.Println("Hello World!")
	redis.InitClient()
	searchAndUpdateMd5()
	// fmt.Printf("\t[Info]New find in files: %v", files)
	sqlite.InitSqlClient()
	sqlite.InsertAdded("test.xyz", "abc.test.xyz")

}
