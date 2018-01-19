package Com

import (
	"shlacd/app/api"
	"fmt"
	"flag"
	"errors"
	"regexp"
	"log"
)

type Get struct{
	Com
}

const usageGet = "usage: \n" +
	"\t  get (\\g) <id>\n" +
	"\t  get (\\g) -c <cron-formatted job line>\n" +
	"\t  get (\\g) --help\n"

func (c *Get) Exec(Tab api.TimeTable, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = c.Usage()
		}

	}(&response, &err)

	var defaultResponse = "null"

	var HELP, HLP, BYCMD bool

	Args := flag.NewFlagSet("com_get", flag.PanicOnError)
	Args.BoolVar(&HELP, "help", false, "show this help")
	Args.BoolVar(&HLP, "h", false, "show this help")
	Args.BoolVar(&BYCMD, "c", false, "find by cron-string")
	Args.Parse(args)

	log.Println(BYCMD, HELP, HLP)


	// show help
	if HELP || HLP || Args.NArg() <1{
		response = c.Usage()

		//show item by index
	}else if BYCMD {

		ws := regexp.MustCompile(`\s+`)
		needle := ws.ReplaceAllString(Args.Arg(0), "")
		for _,j := range Tab.ListJobs(){

			if ws.ReplaceAllString(j.CronLine()+j.Command(),"") == needle{
				response += c.view(j)
			}
		}

	}else{
		if found := Tab.FindJob(Args.Arg(0)); found != nil{
			response += c.view(found)
		}
	}

	if response == "" {
		response = defaultResponse
	}

	return response, nil
}

func (c *Get) Usage() string{

	return c.Desc() + "\n\t" + usageGet
}

func (c *Get) view(job api.Job) string{

	return fmt.Sprintln(
		job.Id(),"\t",
		job.CronLine(),"\t",
		"\""+job.Command()+"\"",
	)
}

