<!--
 * @Author: 1dayluo
 * @Date: 2023-02-07 11:18:40
 * @LastEditTime: 2023-03-11 22:03:55
-->
# SubNya
## <div align="center"><b><a href="README.md">English</a> | <a href="README_CN.md">简体中文</a></b></div>

![](https://img.shields.io/github/commit-activity/w/1dayluo/SubNya_monitor?style=flat-square)    ![](https://img.shields.io/github/license/1dayluo/SubNya_monitor?style=flat-square) 

## Introduction

SubNya_monitor is a new subdomain enumeration and monitoring tool used to track the status of subdomains on the target domain, including newly added and removed subdomains. It utilizes goroutine to increase the speed of subdomain enumeration. The tool stores data in both Redis and SQLite, taking advantage of Redis's speed to monitor changes in file MD5, and SQLite's features to store and update subdomain data, including the use of transactions. Finally, the output result will be saved to a local file record (optional) or sent as a notification to a personal Telegram/email via an API.

The current project has completed its basic functionality, and other functions and the Dockerfile are still under development (see todo below). 

## Easy installation
### release

 you can download them from the [releases](https://github.com/1dayluo/SubNya_monitor/releases/tag/v1.0) page.

### using go install 
If you have a Go environment ready to go (at least go 1.19), it's as easy as:
```lua
go install  github.com/1dayluo/subnya@latest  

```

## Configuration File

The configuration file is located in  `~/.config/subnya/config/config.yml` and looks like this，Executable files need to fill in the absolute path:

```yml
schedule:
  - cron: "* * * *"   # For Dockerfile deployment

monitor:
  dir: 
    - "./test"  # The folder to be monitored - all files in this folder will be traversed
  settings:
    - timeout : 30   # Interval time
    - threads : 10    # Number of threads
    - maxenumerationtime: 10   # Maximum enumeration time
    - outfile : "/var/tmp/"    # Output folder
  
  
redis:
  addr: "172.17.0.1:6379"   # Redis connection address
  password: ""   # Redis password
  db:   0 

sqlite:
  db_1: "./db/monitor.db"   # Set the database storage location
```

## Usage

```lua
Options:
  --update, -u           Check if the monitor has any updates
  --run, -r              Start subdomain finder and update data (including response status code) in SQLite
  --output OUTPUT
  --help, -h             Display this help and exit
```







