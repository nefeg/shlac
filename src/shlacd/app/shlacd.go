package app

import (
	"log"
	"os"
	"time"
	"fmt"
	"shlacd/executor"
	"shlacd/storage"
	"shlacd/client"
	"shlacd/app/Tab"
	. "shared/config/app"
	. "shlacd/app/api"
	"github.com/umbrella-evgeny-nefedkin/slog"
)


type app struct{

	//Api API
	Conf Config

	Tab TimeTable
	Exe Executor

	Client client.Handler
}


func New(AppConfig *Config) *app {

	application := &app{}
	application.Tab     = TimeTable( Tab.New( storage.Resolve(AppConfig.Storage) ))
	application.Exe     = executor.Resolve(AppConfig.Executor)
	application.Client  = client.Resolve(AppConfig.Client)

	return application
}


func (app *app) Run(){

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


	go app.runHrend() // todo remove old jobs

	app.Client.Handle(app.Tab)
}

func (app *app) Stop(code int, message interface{}){

	app.Tab.Close()

	log.Printf("*** Application terminated with message: %s\n\n", message)

	os.Exit(code)
}

func (app *app) runHrend(){

	slog.DebugLn("[app->core] fork")

	for{

		slog.DebugLn("\n**** [app->core] new loop: ", time.Now().String())

		var timeout time.Duration = 60

		if found := app.Tab.ListJobs(); len(found)>0{

			go func(jobs []Job){

				for _, job := range jobs{

					FromTime := time.Now().Add(-1*time.Minute)
					JTS := job.TimeStart( FromTime )

					slog.DebugLn("-------------", job.CronLine(), job.Id())
					slog.DebugLn("now", time.Now().String())
					slog.DebugLn("now-1 (fromTime)", FromTime.String())
					slog.DebugLn("must run", JTS.String())
					slog.DebugLn("diff", time.Since(JTS))

					timeInterval := time.Since(JTS).Seconds()
					if timeInterval >0 && timeInterval <60{
						log.Println("[app] Pulling job:", job.Id())
						if j := app.Tab.PullJob(job.Id()); j != nil{

							log.Println("[app] Job started:", j.Id())
							app.Exe.Exec(job)
							app.Tab.PushJob(job)

						}else{
							log.Println("[app] Pulling job: skip (Can't pull job)", job.Id())

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
