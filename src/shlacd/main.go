package main

import (
	"runtime"
	"shared/sig"
	"os"
	. "shared/config/app"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"github.com/urfave/cli"
	"fmt"
	"shlacd/app"
	"shlacd/app/api"
	"shlacd/app/api/Table"
	"shlacd/storage"
	"shlacd/executor"
	"time"
	"shlacd/client"
	"shlacd/cli/Context"
)

var App Application

var sigIntHandler   = func(){}
var logPrefix       = "[main]"
var ConfigPaths     = []string{
	"config.json",
	"/etc/shlanc/config.json",
	"/etc/shlancd/config.json",
}

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())

	sig.SIG_INT(&sigIntHandler)

	slog.SetLevel(slog.LvlInfo)
	slog.SetFormat(slog.FormatTimed)

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "server version=%s {%s}\n", c.App.Version, c.App.Compiled.Format("2006/01/02 15:04:05"))
	}
}


func main(){

	const FL_DEBUG  = "debug"
	const FL_CONFIG = "config"

	var application Application
	var table       api.Table
	var config      *Config

	defer func(a Application){
		message := "OK"
		code := 0

		r := recover()
		if r != nil{
			message = fmt.Sprint(r)
			code = 1

			if a != nil && a.IsDebug(){slog.PanicLn(r)}
		}

		if a != nil{
			a.Stop(code, message)
		}

	}(application)


	Cli := cli.NewApp()
	Cli.Version             = "0.4"
	Cli.Name                = "ShLAC-server"
	Cli.Usage               = "[SH]lac [L]ike [A]s [C]ron"
	Cli.Author              = "Evgeny Nefedkin"
	Cli.Compiled            = time.Now()
	Cli.Email               = "evgeny.nefedkin@umbrella-web.com"
	Cli.EnableBashCompletion= true
	Cli.Description         = "Distributed and concurrency jobs manager"

	//CONFIG
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


	Cli.Before = func(c *cli.Context) error {

		appOptions := api.AppOptions{}

		if c.GlobalBool(FL_DEBUG){
			slog.SetLevel(slog.LvlDebug)
			slog.DebugF("%s Starting...\n", logPrefix)
			slog.DebugLn(logPrefix, " Args:", os.Args)
			appOptions.DebugMode = true
		}

		// Override config
		if confFile := c.GlobalString(FL_CONFIG); confFile != ""{
			ConfigPaths = []string{confFile}
		}
		slog.DebugLn(logPrefix, " Config paths:", ConfigPaths)


		config = LoadConfig(ConfigPaths)

		table = Table.New( storage.Resolve(config.Storage))

		application = app.New(
			table,
			executor.Resolve(config.Executor),
			appOptions,
		)
		//
		sigIntHandler = func(){
			application.Stop(1, sig.ErrSigINT)
		}

		return nil
	}

	Cli.Action = func(c *cli.Context) {

		go application.Run()

		client.Resolve(config.Client).Handle( Context.New(table) )
	}

	Cli.Run(os.Args)
}