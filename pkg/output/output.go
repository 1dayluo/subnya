/*
 * @Author: 1dayluo
 * @Date: 2023-03-01 15:01:31
 * @LastEditTime: 2023-03-10 19:32:12
 */
package output

import (
	"fmt"
	"io"
	"time"

	"github.com/1dayluo/subnya/pkg/readconf"

	"github.com/fatih/color"
)

type ResultOutput struct {
	Added []string
	Deled []string
	All   []string
	Code  map[string]int
}

var OutFile = readconf.ReadSettingsConfig("outfile") + "/" + time.Now().Format("2006-01-02")

const (
	logo = `


    ███████ ██    ██ ██████  ███    ██ ██    ██  █████  ███    ███  ██████  ███    ██ ██ ████████  ██████  ██████      
    ██      ██    ██ ██   ██ ████   ██  ██  ██  ██   ██ ████  ████ ██    ██ ████   ██ ██    ██    ██    ██ ██   ██     
    ███████ ██    ██ ██████  ██ ██  ██   ████   ███████ ██ ████ ██ ██    ██ ██ ██  ██ ██    ██    ██    ██ ██████      
         ██ ██    ██ ██   ██ ██  ██ ██    ██    ██   ██ ██  ██  ██ ██    ██ ██  ██ ██ ██    ██    ██    ██ ██   ██     
    ███████  ██████  ██████  ██   ████    ██    ██   ██ ██      ██  ██████  ██   ████ ██    ██     ██████  ██   ██     
                                                                                                                       
                                                                                                                       

					Subdomain Monitoring
	  		A tool for monitoring subdomain changes on target domain
			  	Author: @1dayluo(https://github.com/1dayluo)
	--------------------------------------------------------------------------------------------------------------
										`
	usage = `
Usage: subnya_monitor [options]

Options:
	--update, -u           Check update in monitor
	--run, -r              start subdomain finder and update data(include response status code) in sqlite
	--output			   output to file defined in default settings
	--help, -h             display this help and exit
`
)

func PrintLogo() {
	magenta := color.New(color.FgHiMagenta).SprintFunc()
	fmt.Println(magenta(logo))
}

func print(out io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(out, format, args...)
}

func OutResult(results map[string]ResultOutput, out io.Writer, isTTY bool) error {
	/**
	 * @description: Output program running results  (out io.Writer 输出位置，可以是 os.Stdout 或者一个文件句柄)
	 * @return {*}
	 */
	for domain, info := range results {
		print(out, "\n[*]Domain: %v", domain)
		for _, added := range info.Added {
			print(out, "\n   [+]Added: %v", added)
		}
		for _, deled := range info.Deled {
			print(out, "\n   [-]Deled: %v", deled)
		}
		if info.Code != nil {
			print(out, "\n   [MORE]Alive Check Result:\n")
			for subd, code := range info.Code {
				if code != -1 {
					if isTTY {
						var codeColor func(...interface{}) string
						if code >= 200 && code < 300 {
							codeColor = color.New(color.FgGreen).SprintFunc()
						} else if code >= 300 && code < 400 {
							codeColor = color.New(color.FgYellow).SprintFunc()
						} else {
							codeColor = color.New(color.FgRed).SprintFunc()
						}
						print(out, "      - %v ( %v )\n", subd, codeColor(code))

					} else {
						print(out, "      - %v ( %v )\n", subd, code)

					}

				}
			}
		}

	}
	return nil
}
