package app

import (
	"log"
	"os"
	"time"
	"fmt"
	. "shlacd/app/api"
	"github.com/umbrella-evgeny-nefedkin/slog"
)


type application struct{

	table    Table
	executor Executor
	options  AppOptions

}


func New(T Table, E Executor, options AppOptions) *application {

	app := &application{}
	app.table       = T
	app.executor    = E
	app.options     = options

	return app
}

func (app *application) IsDebug() bool{

	return app.options.DebugMode
}


func (app *application) Run(){

	defer func(){
		code    := 0
		message := "no message"
		if r:= recover(); r!=nil{
			log.Println(r)
			code = 1
			message = fmt.Sprint(r)
		}

		app.Stop(code, message)
	}()

	slog.DebugLn("[app] Run")


	app.runHrend() // todo remove old jobs
}

func (app *application) Stop(code int, message interface{}){

	app.table.Close()

	log.Printf("*** Application terminated with message: %s\n\n", message)

	os.Exit(code)
}

func (app *application) runHrend(){

	slog.DebugLn("[app->core] fork")

	for{

		slog.DebugLn("\n**** [app->core] new loop: ", time.Now().String())

		var timeout time.Duration = 60

		if found := app.table.ListJobs(); len(found)>0{

			go func(jobs []Job){

				for _, job := range jobs{

					FromTime := time.Now().Add(-1*time.Minute)
					JTS := job.TimeStart( FromTime )

					slog.DebugLn("-------------", job.TimeLine(), job.Index())
					slog.DebugLn("now: ", time.Now().String())
					slog.DebugLn("now-1 (fromTime): ", FromTime.String())
					slog.DebugLn("must run: ", JTS.String())
					slog.DebugLn("diff: ", time.Since(JTS))

					timeInterval := time.Since(JTS).Seconds()
					if timeInterval >0 && timeInterval <60{
						log.Println("[app] Pulling Job:", job.Index())
						if j := app.table.PullJob(job); j != nil{

							log.Println("[app] Job started:", j.Index())
							app.executor.Exec(job)
							app.table.PushJob(job)

						}else{
							log.Println("[app] Pulling Job: skip (Can't pull Job)", job.Index())

						}
					}
				}

			}(found)
		}

		if timeShift := time.Now().Unix() % int64(timeout); timeShift > 1{// if shift more then 2 seconds
			timeout = timeout - time.Duration(timeShift -1)
		}

		time.Sleep(time.Duration(timeout) * time.Second)
	}
}
