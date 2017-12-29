package Com

import (
	"shlacd/hrontabd"
	"fmt"
	"flag"
	"errors"
)

type List struct{
	Com
}

const usageList = "usage: \n" +
	"\t  list (\\l) \n"

func (c *List) Exec(Tab hrontabd.TimeTable, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = c.Usage()
		}

	}(&response, &err)

	var defaultResponse = "null"

	Args := flag.NewFlagSet("com_list", flag.PanicOnError)
	Args.Parse(args)

	// show help
	for _,job := range Tab.ListJobs() {
		response += c.view( job )
	}

	if response == "" {
		response = defaultResponse
	}

	return response, nil
}

func (c *List) Usage() string{

	return c.Desc() + "\n\t" + usageList
}

func (c *List) view(job hrontabd.Job) string{

	return fmt.Sprintln(
		job.Id(),"\t",
		job.CronLine(),"\t",
		"\""+job.Command()+"\"",
	)
}

