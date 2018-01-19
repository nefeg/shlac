package Com

import (
	"shlacd/app/api"
	"fmt"
	"flag"
	"errors"
)

type List struct{
	Com
}

const usageList = "usage: \n" +
	"\t  list (\\l) \n"

func (c *List) Exec(Tab api.TimeTable, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = c.Usage()
		}

	}(&response, &err)


	Args := flag.NewFlagSet("com_list", flag.PanicOnError)
	Args.Parse(args)

	// show help
	for _,job := range Tab.ListJobs() {
		response += c.view( job )
	}

	return response, nil
}

func (c *List) Usage() string{

	return c.Desc() + "\n\t" + usageList
}

func (c *List) view(job api.Job) string{

	return fmt.Sprintln(
		job.Id(),"\t",
		job.CronLine(),"\t",
		job.Command(),
	)
}

