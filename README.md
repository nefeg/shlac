Distributed Concurrency job manager
=============================

# Table of Contents
1. [Features](#features)
2. [INSTALL](#install)
3. [Configuring](#configuring)
4. [Usage](#usage)
5. [Date format](#date-format)
6. [License](#license)

## Features

 - extended cron syntax
 - remote/local job storage (redis/file)
 - distributing job list for concurrency execution
 - decentralized structure
 
 
## INSTALL


### From Ubuntu repository via apt-get

    # add PPA to apt source list:
    sudo add-apt-repository ppa:onm/shlac
    sudo apt-get update
    
    #install package
    sudo apt-get install shlac



### From .deb package

Download package from [Launchpad PPA](https://launchpad.net/~onm/+archive/ubuntu/shlac/+packages) and install:

        sudo dpkg -i path/to/package.deb

### Compile from source (Required 'go' compiler)
   
   
    # install "go"
    sudo apt-get install golang
      
    # get source from ginhub 
    git clone https://github.com/umbrella-evgeny-nefedkin/shlac
    
    # run build script
    cd shlac; ./build.sh
    
    # change permissions
    chmod +x bin/shlancd
    
    # run server
    bin/shlancd -c config.json
    # or run shlancd as daemon
    # bin/shlancd -c config.json >> shlancd.log &2>1 &

    

## Configuring



### Default config:
        {
            "storage": {
                "type": "redis",
                "options": {
                    "network":  "tcp",
                    "address":  "127.0.0.1:6379",
                    "key":      "shlac"
                }
            },
        
            "client": {
                "type": "socket",
                "options":{
                    "network":  "tcp",
                    "address":  "127.0.0.1:6609"
                }
            },
        
            "executor":{
                "type": "local",
                "options":{
                    "silent":   true,
                    "async":    true
                }
            }
        }



### Parameters


**[STORAGE]**


- **`type:redis`** - use redis as job storage

- `options` - connection options for redis
    
	`options.network` - socket type: `tcp|udp|unix`
    
    `options.address` - socket address(TCP/IP: `127.0.0.1:6379`, unix: `"/path/to/redis.sock"`)
    
    `options.key` - prefix for keys in database

Example:
	
    "storage": {
        "type": "redis",
        "options": {
            "network":  "tcp",
            "address":  "127.0.0.1:6379",
            "key":      "shlac"
        }
    }
 

**[CLIENT]**

- **`type:socket`** - connect via socket
- `options` - connection options for socket
    
	`options.network` - socket type: `tcp|udp|unix`
    
    `options.address` - socket address(TCP/IP: `127.0.0.1:6379`, unix: `"/path/to/redis.sock"`)

```
"client": {
    "type": "socket",
    "options":{
        "network":  "tcp",
        "address":  "127.0.0.1:6607"
    }
}
```
**Just use telnet as client!**


**[EXECUTOR]**

Just use default config:
```
	"executor":{
		"type": "local",
		"options":{
			"silent":   true,
			"async":    true
		}
	}
```


## Usage


- using shlanc cli:


    username:~/$ shlac -h
    NAME:
       ShLAC(client) - [SH]lac [L]ike [A]s [C]ron
    
    USAGE:
       shlac [global options] command [command options] [arguments...]
    
    COMMANDS:
         add, a     Add job
         export, x  Export list of jobs
         remove, r  Remove job by index
         purge      Remove all job
         get, g     Get job by id
         import, i  import jobs from cron-formatted file
         help, h    Shows a list of commands or help for one command
    
    GLOBAL OPTIONS:
       --config value, -c value  path to daemon config-file
       --debug                   show debug log
       --help, -h                show help
       --version, -v             print the version



- using telnet:


    username:/shlac-project$ telnet 127.0.0.1 6609
      Trying 127.0.0.1...
      Connected to 127.0.0.1.
      Escape character is '^]'.
      ShLAC terminal connected OK
      type "help" or "\h" for show available commands
      >>_



## Date format


   This utility uses [modified cronexpr-library](https://github.com/umbrella-evgeny-nefedkin/cronexpr) (here the [origin library](https://github.com/gorhill/cronexpr))  

   Supported extended syntax:
    
    ------------------------------------------------------------------------
    Field name     Mandatory?   Allowed values    Allowed special characters
    ----------     ----------   --------------    --------------------------
    Seconds        No           0-59              * / , -
    Minutes        Yes          0-59              * / , -
    Hours          Yes          0-23              * / , -
    Day of month   Yes          1-31              * / , - L W
    Month          Yes          1-12 or JAN-DEC   * / , -
    Day of week    Yes          0-6 or SUN-SAT    * / , - L #
    Year           No           1970–2099         * / , -


   and aliases:
   
    -------------------------------------------------------------------------------------------------
    Entry       Description                                                             Equivalent to
    -------------------------------------------------------------------------------------------------
    @annually   Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *
    @yearly     Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *
    @monthly    Run once a month at midnight in the morning of the first of the month   0 0 0 1 * * *
    @weekly     Run once a week at midnight in the morning of Sunday                    0 0 0 * * 0 *
    @daily      Run once a day at midnight                                              0 0 0 * * * *
    @hourly     Run once an hour at the beginning of the hour                           0 0 * * * * *
    @reboot     Not supported

   For more information about supported syntax see [documentation of parser](https://github.com/umbrella-evgeny-nefedkin/cronexpr) 

## License

- MIT see <https://github.com/umbrella-evgeny-nefedkin/shlac/blob/master/LICENSE>
