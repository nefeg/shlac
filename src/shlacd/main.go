package main

import (
	"runtime"
	"log"
	"shared/sig"
	"os"
	"io/ioutil"
	"encoding/json"
	. "shared/config/app"
	shlacd "shlacd/app"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

var App Application

const USAGE = "*** Usage: shlacd <path/to/config/file> ***"

func init()  {

	runtime.GOMAXPROCS(runtime.NumCPU())

	slog.SetLevel(slog.LvlInfo)
	slog.SetFormat(slog.FormatTimed)

	sig.SIG_INT(func(){
		log.Println("Terminateing application...")
		// panic(sig.ErrSigINT)

		if App != nil{ // check for application is instantiated
			Application.Stop(App, 1, sig.ErrSigINT)
		}
	})
}


func main(){
	slog.DebugLn("Starting...")

	defer func(){

		log.Println(USAGE)

		if r := recover(); r != nil{
			log.Fatal(r)
		}
	}()

	if len(os.Args) < 2{
		panic("Exit: Expected path to config file")
	}

	slog.DebugLn("[main] Loading config")
	AppConfig := &Config{}
	if config, err := ioutil.ReadFile(os.Args[1]); err == nil{
		json.Unmarshal(config, AppConfig)
	}
	slog.DebugLn("[main] Loaded config: ", AppConfig)


	App = Application( shlacd.New(AppConfig) )
	App.Run()
}