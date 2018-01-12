package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	"errors"
	"shared/sig"
	"net"
	"bufio"
	"regexp"
	"strings"
	"io/ioutil"
	"encoding/json"
	"bytes"
	. "shared/config"
)

var ErrCmdArgs          = errors.New("ERR: expected argument")
var ErrNoConfFile       = errors.New("ERR: config file not found")
var ErrConfCorrupted    = errors.New("ERR: invalid config")
var ConfigPaths         = []string{
	"config.json",
	"/etc/shlac/config.json",
	"/etc/shlacd/config.json",
}

type ExportOpt struct {
	ShowId  bool
}

func init()  {
	sig.SIG_INT(nil)
}


func main(){

	defer func(){
		if r := recover(); r != nil{

			fmt.Println(r)

			if r == ErrCmdArgs{
				fmt.Println("See: shlac <command> --help")
			}
		}
	}()



	app := cli.NewApp()
	app.Version             = "0.1"
	app.Name                = "ShLAC"
	app.Usage               = "[SH]lac [L]ike [A]s [C]ron"
	app.Author              = "Evgeny Nefedkin"
	app.Email               = "evgeny.nefedkin@umbrella-web.com"
	app.EnableBashCompletion= true
	app.Description         = "Distributed and concurrency job manager\n" +

		"\t\tSupported extended syntax:\n" +
		"\t\t------------------------------------------------------------------------\n" +
		"\t\tField name     Mandatory?   Allowed values    Allowed special characters\n" +
		"\t\t----------     ----------   --------------    --------------------------\n" +
		"\t\tSeconds        No           0-59              * / , -\n" +
		"\t\tMinutes        Yes          0-59              * / , -\n" +
		"\t\tHours          Yes          0-23              * / , -\n" +
		"\t\tDay of month   Yes          1-31              * / , - L W\n" +
		"\t\tMonth          Yes          1-12 or JAN-DEC   * / , -\n" +
		"\t\tDay of week    Yes          0-6 or SUN-SAT    * / , - L #\n" +
		"\t\tYear           No           1970–2099         * / , -\n" +

		"\n\n" +

		"\t\tand aliases:\n" +
		"\t\t-------------------------------------------------------------------------------------------------\n" +
		"\t\tEntry       Description                                                             Equivalent to\n" +
		"\t\t-------------------------------------------------------------------------------------------------\n" +
		"\t\t@annually   Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *\n" +
		"\t\t@yearly     Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *\n" +
		"\t\t@monthly    Run once a month at midnight in the morning of the first of the month   0 0 0 1 * * *\n" +
		"\t\t@weekly     Run once a week at midnight in the morning of Sunday                    0 0 0 * * 0 *\n" +
		"\t\t@daily      Run once a day at midnight                                              0 0 0 * * * *\n" +
		"\t\t@hourly     Run once an hour at the beginning of the hour                           0 0 * * * * *\n" +
		"\t\t@reboot     Not supported"


	// CONFIG
	app.Flags =  []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to daemon config-file",
		},
	}


	// COMMANDS
	app.Commands = []cli.Command{

		{// REMOVE
			Name:    "remove",
			Aliases: []string{"rm", "r"},
			Usage:   "remove jobs ",
			UsageText: "Example: \n" +
				"\t\tshlac rm -i <job id>\n" +
				"\t\tshlac rm --all",

			Flags: 	[]cli.Flag{
				cli.BoolFlag{
					Name:  "all,purge",
					Usage: "remove all jobs",
				},
				cli.StringFlag{
					Name:  "id,i",
					Usage: "remove job by id",
				},
			},

			Action:  func(c *cli.Context) error {

				// Override config
				if confFile := c.GlobalString("config"); confFile != ""{
					ConfigPaths = []string{confFile}
				}

				connection := connect( loadConfig(ConfigPaths) )
				defer func(){
					connection.Write([]byte(`\q`))
					connection.Close()
				}()

				if jobId := c.String("id"); jobId != ""{
					remove(jobId)

				}else if c.Bool("all"){
					purge()
				}

				return nil
			},
		},

		{// PURGE
			Name:    "purge",
			Usage:   "remove all jobs ",
			UsageText: "Example: " +
				"shlac purge",

			Action:  func(c *cli.Context) error {

				// Override config
				if confFile := c.GlobalString("config"); confFile != ""{
					ConfigPaths = []string{confFile}
				}

				purge()

				return nil
			},
		},

		{// IMPORT
			Name:    "import",
			Aliases: []string{"i"},
			Usage:   "import jobs from cron-formatted file",
			UsageText: "Example: " +
				"shlac import <path/to/import/file>",

			Flags: 	[]cli.Flag{
				cli.BoolFlag{
					Name:  "purge",
					Usage: "delete jobs before import",
				},

				cli.BoolFlag{
					Name:  "skip-check, s",
					Usage: "add job even if same is already exist (skip checking for duplicates)",
				},
			},



			Action:  func(c *cli.Context) error {

				filePath := c.Args().Get(0)
				if filePath == "" {
					panic(ErrCmdArgs)
				}

				// Override config
				if confFile := c.GlobalString("config"); confFile != ""{
					ConfigPaths = []string{confFile}
				}

				// clean table before import
				if c.Bool("purge"){ purge() }

				Import(filePath, !c.Bool("skip-check"))

				return nil
			},
		},

		{// ADD JOB
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add job from cron-formatted line",
			UsageText: "Example: " +
				"shlac add '<cron-formatted line>'",

			Flags: 	[]cli.Flag{
				cli.BoolFlag{
					Name:  "skip-check, s",
					Usage: "add job even if same is already exist (skip checking for duplicates)",
				},
			},

			Action:  func(c *cli.Context) error {

				cronString := c.Args().Get(0)
				if cronString == "" {
					panic(ErrCmdArgs)
				}

				// Override config
				if confFile := c.GlobalString("config"); confFile != ""{
					ConfigPaths = []string{confFile}
				}

				ImportLine(cronString, !c.Bool("skip-check"))

				return nil
			},
		},

		{// EXPORT
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "export jobs to file in cron-format",
			UsageText: "Example: \n" +
				"\t\tto stdout:\tshlac export\n" +
				"\t\tto file:\tshlac export -f <path/to/export/file>",
			Flags: 	[]cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "export to file",
				},
				cli.BoolFlag{
					Name:   "show-id, i",
					Usage:  "export with job ids",
				},
			},
			Action:  func(c *cli.Context) error {

				// Override config
				if confFile := c.GlobalString("config"); confFile != ""{
					ConfigPaths = []string{confFile}
				}

				exportOptions   := ExportOpt{ShowId:c.Bool("show-id")}
				exportedData    := Export(exportOptions)


				if exportFile := c.String("file"); exportFile != ""{
					ioutil.WriteFile(exportFile, []byte(exportedData), 0644)

				}else{
					fmt.Print(exportedData)
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}


func Export(options ExportOpt) string{

	response := sendCommand(`\l`) // send command and read answer

	response = response[:len(response)-4] // remove terminal bytes

	if !options.ShowId{
		re := regexp.MustCompile(`(?m)^.+?\s+`) // remove jobs id
		response = re.ReplaceAll(response, []byte{})

	}else{

		matches := regexp.MustCompile(`(?m)^(.+?\s)([0-9,\/*-LW#]+\s+){5,7}(.+)$`).FindAllSubmatchIndex(response, -1)

		for i := range matches{

			response = append(response, response[matches[i][3]:matches[i][7]]...)
			response = append(response, []byte(" # ")...)
			response = append(response, response[:matches[i][3]]...)
		}
	}

	return string(response)
}

func Import(filePath string, checkDuplicates bool){

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner     := bufio.NewScanner(file)

	for scanner.Scan() {

		ImportLine(scanner.Text(), checkDuplicates)

		// MAIN LOOP
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}

func ImportLine(cronString string, checkDuplicates bool){

	var importLine string

	if cronString == "\n"{
		importLine = cronString

	}else{
		re := regexp.MustCompile(`^([0-9,\/*-LW#]+\s+){5,7}(.+)$`)
		matches := re.FindStringSubmatchIndex(cronString)

		cronLine    := strings.Trim( cronString[:matches[4]], " \t" )
		commandLine := strings.Trim( cronString[matches[4]:], " \t" )

		importLine = fmt.Sprintf(`\a -cron "%s" -cmd %q`, cronLine, commandLine)

		if cronLine[:1] == `#`{
			fmt.Printf("SKIPP (disabled)>> %s\n", importLine)
			return
		}

		if checkDuplicates && isDuplicated(cronString){
			fmt.Printf("SKIPP (duplicated)>> %s\n", importLine)
			return
		}


		fmt.Printf("IMPORT>> %s\n", importLine)
	}


	sendCommand(importLine)
}



func loadConfig(configPaths []string) (config *Config) {

	configRaw := func(configPaths []string) (configRaw []byte){

		for _,configPath := range configPaths{

			configRaw, err := ioutil.ReadFile(configPath)

			if err == nil && configRaw != nil {
				return configRaw
			}
		}

		return nil

	}(configPaths)


	if configRaw == nil {
		panic(fmt.Sprint(ErrNoConfFile, configPaths))
	}

	config = &Config{}
	if err := json.Unmarshal(configRaw, config); err != nil{
		panic(ErrConfCorrupted)
	}

	return config
}

func connect(config *Config) (connection net.Conn){
	if config.Client.Type != "socket" {
		panic("Unsupported client type")
	}

	conn, err := net.Dial(config.Client.Options.Network, config.Client.Options.Address)
	if err != nil{
		panic(err)
	}

	return conn
}

func sendCommand(command string) (received []byte){

	connection := connect( loadConfig(ConfigPaths) )
	defer func(){
		connection.Write([]byte(`\q`))
		connection.Close()
	}()

	flushConnection(connection) // clear socket buffer

	connection.Write([]byte(command+"\n"))

	received = flushConnection(connection)

	return received
}

func flushConnection(connection net.Conn) (flushed []byte){

	bufSize := 256
	buf := make([]byte, bufSize)

	for{
		n,e := connection.Read(buf)

		flushed = append(flushed, buf[:n]...)

		if e != nil || n < bufSize {break}
	}

	return flushed
}

func purge(){
	sendCommand(`\r --all`)
}

func remove(jobId string){
	sendCommand(`\r -id `+jobId)
}

func isDuplicated(cronLine string) bool {

	connection := connect( loadConfig(ConfigPaths) )
	defer func(){
		connection.Write([]byte(`\q`))
		connection.Close()
	}()

	cronLine = strings.Replace(cronLine, `"`, `\"`, -1)

	flushConnection(connection)
	connection.Write([]byte(`\g -c "` +cronLine+ `"`))

	response := make([]byte, 8)
	connection.Read(response)

	return !bytes.Equal(response, []byte{110, 117, 108, 108, 0, 10, 62, 62})
}