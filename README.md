# SubNya_monitor
## <div align="center"><b><a href="README.md">English</a> | <a href="README_CN.md">简体中文</a></b></div>

## Introduction

SubNya_monitor is a new subdomain enumeration and monitoring tool used to track the status of subdomains on the target domain, including newly added and removed subdomains. It utilizes goroutine to increase the speed of subdomain enumeration. The tool stores data in both Redis and SQLite, taking advantage of Redis's speed to monitor changes in file MD5, and SQLite's features to store and update subdomain data, including the use of transactions. Finally, the output result will be saved to a local file record (optional) or sent as a notification to a personal Telegram/email via an API.

The current project has completed its basic functionality, and other functions and the Dockerfile are still under development (see todo below). The original plan included the following:

- Provide interfaces and demos of different instant messaging tools to notify users of newly discovered subdomains.
- Set up scheduled tasks to periodically query new subdomains.
- Support periodic monitoring of folders and command line input.
- Newly added subdomains will be sent to the message queue, and users can customize which applications they want to be notified.
- Built-in subdomain collection/extension using subfinder to query subdomains.

## Configuration File

The configuration file is located in `./config/config.yml` and looks like this:

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



## TODO:

1. Implement the output traversal result function and confirm the output format.
2. Optimize error handling and record errors in logs.
3. Design notification function to ensure that newly added subdomains can be notified to users.
4. Beautify console output.
5. Set up scheduled tasks and corresponding configuration files.
6. Optimize threading, use proxies or other details, etc.
7. English Readme document.



