package main

import (
	"DomainMonitor/pkg/db"
	"DomainMonitor/pkg/readconf"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
	db.InitClient()
	readconf.ReadMonitorDir()

}
