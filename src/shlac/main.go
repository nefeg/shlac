package main

import (
	"os"
	"fmt"
	"errors"
	"shared/sig"
	shared "shared/config/app"
	"github.com/umbrella-evgeny-nefedkin/slog"

	capi "shlacd/cli"
	capiCtx "shlacd/cli/Context"
	"shlacd/app/api/Table"
	"github.com/urfave/cli"
	"shlacd/storage"
)

var ConfigPaths         = []string{
	"config.json",
	"/etc/shlac2/config.json",
	"/etc/shlacd/config.json",
}


var ErrCmdArgs          = errors.New("ERR: expected argument")
var sigIntHandler       = func(){} // check

const FL_DEBUG      = "debug"

func init()  {


	sig.SIG_INT(&sigIntHandler)

	slog.SetLevel(slog.LvlCrit)
	slog.SetFormat(slog.FormatTimed)
}


func main(){

	var CliContext capi.Context

	sigIntHandler = func(){
		CliContext.Term()
	}

	defer func(a *capi.Context){
		if r := recover(); r != nil{

			slog.DebugLn("[main->defer] ", r)

			if r == ErrCmdArgs || r == shared.ErrNoConfFile{
				fmt.Println(r)
				fmt.Println("See: `shlac --help` or `shlac <command> --help`")

			}else if r == shared.ErrConfCorrupted{
				fmt.Println(r)

			}else{
				fmt.Println(r)
			}
		}
		if *a != nil{
			(*a).Term()
		}

	}(&CliContext)


	Cli := capi.New()

	Cli.Before = func(c *cli.Context) error {

		// debug flag
		if c.GlobalBool(FL_DEBUG){
			slog.SetLevel(slog.LvlDebug)
			slog.DebugLn("[main] os.Args: ", os.Args)
		}

		// Override config
		if confFile := c.GlobalString("config"); confFile != ""{
			ConfigPaths = []string{confFile}
		}
		slog.DebugLn("[main ->Cli.Before] Config paths: ", ConfigPaths)



		mainConfig := shared.LoadConfig(ConfigPaths)

		JTable := Table.New( storage.Resolve(mainConfig.Storage) )
		slog.DebugLn("[main ->Cli.Before] JTable: ", JTable)

		CliContext = capiCtx.New( JTable )
		slog.DebugLn("[main ->Cli.Before] CliContext: ", CliContext)

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
		capi.NewComAdd(&CliContext),
		capi.NewComExport(&CliContext),
		capi.NewComRemove(&CliContext),
		capi.NewComPurge(&CliContext),
		capi.NewComGet(&CliContext),
		capi.NewComImport(&CliContext),
	}

	Cli.Action = func(c *cli.Context) error {

		c.App.Writer.Write([]byte("shlac(client) version: "+c.App.Version+"\nSee `shlac --help` for more information\n"))

		return nil
	}

	Cli.Run(os.Args)
}




//func main(){
//
//	defer func(){
//		if r := recover(); r != nil{
//
//			fmt.Println(r)
//
//			if r == ErrCmdArgs{
//				fmt.Println("See: shlac2 <command> --help")
//			}
//		}
//	}()
//
//
//	Cli := capi.New()
//	Cli.Before = func(c *cli.Context) error {
//		// Override config
//		if confFile := c.GlobalString("config"); confFile != ""{
//			ConfigPaths = []string{confFile}
//		}
//
//		if c.GlobalBool(FL_DEBUG){
//			slog.SetLevel(slog.LvlDebug)
//		}
//
//		slog.DebugLn("[main->Cli.Before] Config paths:", ConfigPaths)
//
//		AppConfig := app.LoadConfig(ConfigPaths)
//
//
//		ComSender := lib.NewCommander()
//		ComSender.Connect( &addr.Config{
//			Protocol: AppConfig.Client.Options.Network,
//			Address: AppConfig.Client.Options.Address,
//		})
//
//		Application = shlac2.New( ComSender )
//
//		return nil
//	}
//
//
//	// CONFIG
//	Cli.Flags =  []cli.Flag{
//		cli.StringFlag{
//			Name:  "config, c",
//			Usage: "path to daemon config-file",
//		},
//		cli.BoolFlag{
//			Name:  "debug",
//			Usage: "show debug log",
//		},
//	}
//
//	Cli.Commands = []cli.Command{
//		capi.NewComGet(con)
//	}
//
//	// COMMANDS
//	//Cli.Commands = []cli.Command{
//	//
//	//	{// REMOVE
//	//		Name:    COM_REMOVE,
//	//		Aliases: []string{"rm", "r"},
//	//		Usage:   "remove jobs ",
//	//		UsageText: "Example: \n" +
//	//			"\t\tshlac rm <Job id>\n" +
//	//			"\t\tshlac rm --all",
//	//
//	//		Flags: 	[]cli.Flag{
//	//			cli.BoolFlag{
//	//				Name:  "all,purge",
//	//				Usage: "remove all jobs",
//	//			},
//	//		},
//	//
//	//		Action:  func(c *cli.Context) error {
//	//
//	//			if c.Bool("all"){
//	//				Application.Purge()
//	//
//	//			}else if jobId := c.Args().Get(0); jobId != ""{
//	//				Application.Remove(jobId)
//	//
//	//			}else{
//	//				panic(ErrCmdArgs)
//	//			}
//	//
//	//			return nil
//	//		},
//	//	},
//	//
//	//	{// PURGE
//	//		Name:    COM_PURGE,
//	//		Usage:   "remove all jobs ",
//	//		UsageText: "Example: " +
//	//			"shlac2 purge",
//	//
//	//		Action:  func(c *cli.Context) error {
//	//
//	//			Application.Purge()
//	//
//	//			return nil
//	//		},
//	//	},
//	//
//	//	{// IMPORT
//	//		Name:    COM_IMPORT,
//	//		Aliases: []string{"i"},
//	//		Usage:   "import jobs from cron-formatted file",
//	//		UsageText: "Example: " +
//	//			"shlac2 import <path/to/import/file>",
//	//
//	//		Flags: 	[]cli.Flag{
//	//			cli.BoolFlag{
//	//				Name:  "purge",
//	//				Usage: "delete jobs before import",
//	//			},
//	//
//	//			cli.BoolFlag{
//	//				Name:  "skip-check, s",
//	//				Usage: "add Job even if same is already exist (skip checking for duplicates)",
//	//			},
//	//		},
//	//
//	//
//	//
//	//		Action:  func(c *cli.Context) error {
//	//
//	//			filePath := c.Args().Get(0)
//	//			if filePath == "" {
//	//				panic(ErrCmdArgs)
//	//			}
//	//
//	//			// clean Table before import
//	//			if c.Bool("purge"){ Application.Purge() }
//	//
//	//			Application.Import(filePath, !c.Bool("skip-check"))
//	//
//	//			return nil
//	//		},
//	//	},
//	//
//	//	{// ADD JOB
//	//		Name:    COM_ADD,
//	//		Aliases: []string{"a"},
//	//		Usage:   "add Job from cron-formatted line",
//	//		UsageText: "Example: " +
//	//			"shlac2 add '<cron-formatted line>'",
//	//
//	//		Flags: 	[]cli.Flag{
//	//			cli.BoolFlag{
//	//				Name:  "skip-check, s",
//	//				Usage: "add Job even if same is already exist (skip checking for duplicates)",
//	//			},
//	//		},
//	//
//	//		Action:  func(c *cli.Context) error {
//	//
//	//			cronString := c.Args().Get(0)
//	//			if cronString == "" {
//	//				panic(ErrCmdArgs)
//	//			}
//	//
//	//			Application.ImportLine(cronString, !c.Bool("skip-check"))
//	//
//	//			return nil
//	//		},
//	//	},
//	//
//	//	{// EXPORT
//	//		Name:    COM_EXPORT,
//	//		Aliases: []string{"e"},
//	//		Usage:   "export jobs to file in cron-format",
//	//		UsageText: "Example: \n" +
//	//			"\t\tto stdout:\tshlac export\n" +
//	//			"\t\tto file:\tshlac export -f <path/to/export/file>",
//	//		Flags: 	[]cli.Flag{
//	//			cli.StringFlag{
//	//				Name:  "file, f",
//	//				Usage: "export to file",
//	//			},
//	//			cli.BoolFlag{
//	//				Name:   "show-id, i",
//	//				Usage:  "export with Job ids",
//	//			},
//	//		},
//	//		Action:  func(c *cli.Context) error {
//	//
//	//			exportOptions   := shlac2.ExportOpt{ShowId:c.Bool("show-id")}
//	//			exportedData    := Application.Export(exportOptions)
//	//
//	//
//	//			if exportFile := c.String("file"); exportFile != ""{
//	//				ioutil.WriteFile(exportFile, []byte(exportedData), 0644)
//	//
//	//			}else{
//	//				fmt.Print(exportedData)
//	//			}
//	//
//	//			return nil
//	//		},
//	//	},
//	//}
//
//	Cli.Run(os.Args)
//}


