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
)

var App Application

const USAGE = "*** Usage: shlacd <path/to/config/file> ***"

func init()  {

	runtime.GOMAXPROCS(runtime.NumCPU())

	sig.SIG_INT(func(){
		log.Println("Terminateing application...")
		// panic(sig.ErrSigINT)

		Application.Stop(App, 1, sig.ErrSigINT)
	})
}


func main(){
	log.Println("Starting...")

	defer func(){

		log.Println(USAGE)

		if r := recover(); r != nil{
			log.Fatal(r)
		}
	}()

	if len(os.Args) < 2{
		panic("Exit: Expected path to config file")
	}

	AppConfig := Config{}
	if config, err := ioutil.ReadFile(os.Args[1]); err == nil{
		json.Unmarshal(config, &AppConfig)
	}

	App = Application( shlacd.New(AppConfig) )
	App.Run()
}