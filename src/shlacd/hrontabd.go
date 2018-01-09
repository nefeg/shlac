package main

import (
	"log"
	"os"
	"time"
	"fmt"
	"shlacd/hrontabd/Tab"
	"shlacd/hrontabd"
	"shlacd/executor"
	"shlacd/storage"
	"shlacd/client"
	. "shared/config"
)


type app struct{

	//Api API
	Conf Config

	Tab hrontabd.TimeTable
	Exe hrontabd.Executor

	Client client.Handler
}


func CreateApp(AppConfig Config) *app {

	application := &app{}
	application.Tab     = hrontabd.TimeTable( Tab.New( storage.Resolve(AppConfig.Storage) ))
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

	go app.runHrend() // todo remove old jobs

	app.Client.Handle(app.Tab)
}

func (app *app) Stop(code int, message interface{}){

	app.Tab.Close()

	log.Printf("*** Application terminated with message: %s\n\n", message)

	os.Exit(code)
}

func (app *app) runHrend(){

	for{
		var timeout time.Duration = 60

		if found := app.Tab.ListJobs(); len(found)>0{

			go func(jobs []hrontabd.Job){

				for _, job := range jobs{

					JTS := job.TimeStart( time.Now().Add(-time.Duration(timeout)*time.Second) )

					//log.Println("-------------", job.CronLine())
					//log.Println("Job", job.Id(), "now", time.Now().String())
					//log.Println("Job", job.Id(), "jts",JTS.String())
					//log.Println("Job", job.Id(), "since", time.Since(JTS))
					//log.Println("-------------")

					timeInterval := time.Since(JTS).Seconds()
					if timeInterval >0{
						log.Println("[hrontabd] Pulling job:", job.Id())
						if j := app.Tab.PullJob(job.Id()); j != nil{

							log.Println("[hrontabd] Job started:", j.Id())
							app.Exe.Exec(job)
							app.Tab.PushJob(job)

						}else{
							log.Println("[hrontabd] Pulling job: skip (Can't pull job)", job.Id())

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