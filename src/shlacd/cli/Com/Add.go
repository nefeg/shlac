package Com

import (
	"shlacd/hrontabd"
	"fmt"
	"flag"
	"shlacd/hrontabd/Job"
	"errors"
	"log"
	"github.com/satori/go.uuid"
)

type Add struct{
	Com
}

const usageAdd =
	"\\a -cron <cron-line> -cmd <command> [-id <id>] [--force]\n" +
	"\t\\a --help\n"


func (c *Add)Exec(Tab hrontabd.TimeTable, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = c.Usage()

			log.Println("[ComAdd]Exec: ", r)
		}

	}(&response, &err)

	var INDEX, CMD, CLINE string
	var OVERRIDE, HELP, HLP bool

	Args := flag.NewFlagSet("com_add", flag.PanicOnError)
	Args.StringVar(&INDEX,      "id",       "",     "record index(name/id)? unique string")
	Args.StringVar(&CLINE,      "cron",     "",     "cron-formatted time line")
	Args.StringVar(&CMD,        "cmd",      "",     "command")
	Args.BoolVar(&OVERRIDE,     "force",    false,  "allow to override existed records")
	Args.BoolVar(&HELP,         "help",     false,  "show this help")
	Args.BoolVar(&HLP,          "h",        false,  "show this help")
	Args.Parse(args)


	if HELP || HLP || CMD=="" || CLINE == ""{
		response = c.Usage()

	}else{

		if INDEX==""{
			if uid, err := uuid.NewV4(); err == nil{
				INDEX = uid.String()
			}else{
				panic(err)
			}
		}

		job := Job.New()
		job.SetID(INDEX)
		job.SetCronLine(CLINE)
		job.SetCommand(CMD)

		log.Println(job)

		Tab.AddJob(job, OVERRIDE)

		response = "OK"
	}


	return response, err
}

func (c *Add) Usage() string{
	return c.Desc() + "\n\t" + usageAdd
}