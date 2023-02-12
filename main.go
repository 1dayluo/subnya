package main

import (
	"DomainMonitor/pkg"
	"DomainMonitor/pkg/db"
	"DomainMonitor/pkg/readconf"
	"fmt"
)

func InsertNewFindMd5(fname string, fmd5 string) {
	db.SetMd5InDB(fmd5)
	db.UpdateFileMd5(fname, fmd5)

}
func searchAndUpdateMd5() (newMonitorFiles []string) {
	//@title searchAndUpdateMd5
	//@param
	//Return Nil
	dirs := readconf.ReadMonitorDir()
	var dirInfo []map[string]string
	for _, dir := range dirs {
		dirInfo = append(dirInfo, pkg.ReadFromDir(dir))
	}
	for _, finfos := range dirInfo {
		for fname, fmd5 := range finfos {
			if db.CheckMd5InDB(fmd5) {
				// fmt.Println("\t[Info]Exists:", fname)
				history_md5 := db.SearchFileMd5(fname)
				if history_md5 != fmd5 {
					InsertNewFindMd5(fname, fmd5)
					newMonitorFiles = append(newMonitorFiles, fname)
				}
			} else { //If fmd5 not in
				// fmt.Println("\t[Info]Not Exists:", fname)
				db.UpdateFileMd5(fname, fmd5)
				InsertNewFindMd5(fname, fmd5)
				newMonitorFiles = append(newMonitorFiles, fname)
			}
		}
	}
	return
}
func main() {
	fmt.Println("Hello World!")
	db.InitClient()
	files := searchAndUpdateMd5()
	fmt.Printf("\t[Info]New find in files: %v", files)

}
