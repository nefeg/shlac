package Com

import (
	"shlacd/app/api"
	"fmt"
	"flag"
	"errors"
	"log"
)

type Remove struct{
	Com
}

const usageRm = "usage: \n" +
	"\t  rm (\\r) -id <index> \n" +
	"\t  rm (\\r) --all \n" +
	"\t  rm (\\r) --help\n"

func (c *Remove)Exec(Tab api.TimeTable, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = c.Usage()
		}

	}(&response, &err)

	var ID string
	var ALL, HELP bool

	Args := flag.NewFlagSet("com_remove", flag.PanicOnError)
	Args.StringVar(&ID, "id", "", "remove record by id")
	Args.BoolVar(&ALL, "all", false, "remove all records")
	Args.BoolVar(&HELP, "help", false, "show this help")
	Args.Parse(args)

	response = c.Usage()

	if HELP || Args.NFlag() <1 {
		response = c.Usage()

	}else if ALL{
		Tab.Flush()
		response = "OK"

	}else if ID != ""{

		if Tab.RmJob(ID) {
			response = "OK"

		}else{
			if Tab.FindJob(ID) == nil{
				response = "job not found"
			}else{
				response = "can`t remove job (unknown error)"
				log.Println("[com.remove] Exec:", response)
			}
		}
	}

	return response, err
}

func (c *Remove) Usage() string{
	return c.Desc() + "\n\t" + usageRm
}