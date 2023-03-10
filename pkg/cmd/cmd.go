/*
 * @Author: 1dayluo
 * @Date: 2023-02-19 14:03:48
 * @LastEditTime: 2023-03-10 19:30:36
 */
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/1dayluo/subnya/pkg/readconf"

	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

var threads, _ = strconv.Atoi(readconf.ReadSettingsConfig("threads"))
var timeout, _ = strconv.Atoi(readconf.ReadSettingsConfig("timeout"))
var maxenumerationtime, _ = strconv.Atoi(readconf.ReadSettingsConfig("maxenumerationtime"))

func FindStart(domain string) string {
	//@title FindStart
	//@param
	//Return
	runnerInstance, err := runner.NewRunner(&runner.Options{
		Threads:            threads,                  // Thread controls the number of threads to use for active enumerations
		Timeout:            timeout,                  // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: maxenumerationtime,       // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		Resolvers:          resolve.DefaultResolvers, // Use the default list of resolvers by marshaling it to the config
		ResultCallback: func(s *resolve.HostEntry) { // Callback function to execute for available host
			log.Println(s.Host, s.Source)
		},
	})

	buf := bytes.Buffer{}
	err = runnerInstance.EnumerateSingleDomain(domain, []io.Writer{&buf})
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(&buf)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("data: %s", data)

	return fmt.Sprintf("%s", data)

}
