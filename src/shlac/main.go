package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	"errors"
	"io/ioutil"
	"encoding/json"
	"time"
	"shared/sig"
	"shared/config/app"
	"shared/config/addr"
	"github.com/umbrella-evgeny-nefedkin/slog"

	"shlac/lib"
	shlac "shlac/app"
)

var Application App

var ConfigPaths         = []string{
	"config.json",
	"/etc/shlac/config.json",
	"/etc/shlacd/config.json",
}


var ErrCmdArgs          = errors.New("ERR: expected argument")
var ErrNoConfFile       = errors.New("ERR: config file not found")
var ErrConfCorrupted    = errors.New("ERR: invalid config data")

const FL_DEBUG      = "debug"

const COM_IMPORT    = "import"
const COM_ADD       = "add"
const COM_EXPORT    = "export"
const COM_PURGE     = "purge"
const COM_REMOVE    = "remove"




func init()  {
	sig.SIG_INT(nil)

	slog.SetLevel(slog.LvlError)
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


	Cli := cli.NewApp()
	Cli.Version             = "0.1"
	Cli.Name                = "ShLAC"
	Cli.Usage               = "[SH]lac [L]ike [A]s [C]ron"
	Cli.Author              = "Evgeny Nefedkin"
	Cli.Compiled            = time.Now()
	Cli.Email               = "evgeny.nefedkin@umbrella-web.com"
	Cli.EnableBashCompletion= true
	Cli.Description         = "Distributed and concurrency job manager\n" +

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


	Cli.Before = func(c *cli.Context) error {
		// Override config
		if confFile := c.GlobalString("config"); confFile != ""{
			ConfigPaths = []string{confFile}
		}

		if c.GlobalBool(FL_DEBUG){
			slog.SetLevel(slog.LvlDebug)
		}

		slog.DebugLn("Config paths:", ConfigPaths)

		AppConfig := loadConfig(ConfigPaths)

		ComSender := lib.NewCommander()
		ComSender.Connect( &addr.Config{
			Protocol: AppConfig.Client.Options.Network,
			Address: AppConfig.Client.Options.Address,
		})

		Application = shlac.New( ComSender )

		return nil
	}


	// CONFIG
	Cli.Flags =  []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to daemon config-file",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "show debug log",
		},
	}


	// COMMANDS
	Cli.Commands = []cli.Command{

		{// REMOVE
			Name:    COM_REMOVE,
			Aliases: []string{"rm", "r"},
			Usage:   "remove jobs ",
			UsageText: "Example: \n" +
				"\t\tshlac rm <job id>\n" +
				"\t\tshlac rm --all",

			Flags: 	[]cli.Flag{
				cli.BoolFlag{
					Name:  "all,purge",
					Usage: "remove all jobs",
				},
			},

			Action:  func(c *cli.Context) error {

				if c.Bool("all"){
					Application.Purge()

				}else if jobId := c.Args().Get(0); jobId != ""{
					Application.Remove(jobId)

				}else{
					panic(ErrCmdArgs)
				}

				return nil
			},
		},

		{// PURGE
			Name:    COM_PURGE,
			Usage:   "remove all jobs ",
			UsageText: "Example: " +
				"shlac purge",

			Action:  func(c *cli.Context) error {

				Application.Purge()

				return nil
			},
		},

		{// IMPORT
			Name:    COM_IMPORT,
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

				// clean table before import
				if c.Bool("purge"){ Application.Purge() }

				Application.Import(filePath, !c.Bool("skip-check"))

				return nil
			},
		},

		{// ADD JOB
			Name:    COM_ADD,
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

				Application.ImportLine(cronString, !c.Bool("skip-check"))

				return nil
			},
		},

		{// EXPORT
			Name:    COM_EXPORT,
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

				exportOptions   := shlac.ExportOpt{ShowId:c.Bool("show-id")}
				exportedData    := Application.Export(exportOptions)


				if exportFile := c.String("file"); exportFile != ""{
					ioutil.WriteFile(exportFile, []byte(exportedData), 0644)

				}else{
					fmt.Print(exportedData)
				}

				return nil
			},
		},
	}

	Cli.Run(os.Args)
}


func loadConfig(configPaths []string) *app.Config{

	configRaw := func(configPaths []string) (configRaw []byte){

		for _,configPath := range configPaths{

			configRaw, err := ioutil.ReadFile(configPath)

			if err == nil && configRaw != nil {
				slog.DebugLn("Loaded config:", configPath)
				return configRaw
			}
		}

		return nil

	}(configPaths)


	if configRaw == nil {
		panic(fmt.Sprint(ErrNoConfFile, configPaths))
	}

	config := &app.Config{}
	if err := json.Unmarshal(configRaw, config); err != nil{
		panic(fmt.Sprint(ErrConfCorrupted, err))
	}

	return config
}



