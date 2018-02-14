package cli

import (
	"shlacd/app/api"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

type viewOptions struct{

	ShowIndex   bool
	CronLegacy  bool
}

func view(jobs []api.Job, options viewOptions) (response string){

	for _,job := range jobs{
		response += viewItem(job, options)
	}

	return response
}

func viewItem(job api.Job, options viewOptions) (response string){

	slog.Infoln("[libCom->viewItem] ", job.String(), options)

	response = job.String()

	if options.ShowIndex{
		if options.CronLegacy{
			response += " # " + job.Index()
		}else{
			response = job.Index() + " " + response

		}
	}

	response += "\n"

	return response
}

