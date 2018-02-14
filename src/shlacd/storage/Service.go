package storage

import (
	"shlacd/app/api"
	"shlacd/app/api/Job"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"github.com/satori/go.uuid"
)

func New(adapter Adapter) api.Storage{
	return &service{adapter}
}


// must implement api.Storage interface
type service struct{

	adapter Adapter
}

func (s service)Add(job api.Job, force bool) (api.Job, error){

	// generate index if not given
	if job.Index() == ""{
		if uid, err := uuid.NewV4(); err == nil{
			job.IndexX( uid.String() )
			slog.Debugln("[storage.redis -> Add] Generated index: ", job.Index())

		}else{
			panic(err)
		}
	}


	if !s.adapter.Lock(job.Index()){
		panic("can't get lock")
	}
	defer s.adapter.UnLock(job.Index())


	return job, s.adapter.Add(job.Index(), job.Serialize(), force)
}


func (s service)Load(job api.Job) (api.Job, error){

	record := s.adapter.Get(job.Index())

	return job, job.UnSerialize(record)
}


func (s service)Rm(job api.Job) bool{
	return s.adapter.Rm(job.Index())
}


func (s service)List() []api.Job{

	var jobList []api.Job

	if records := s.adapter.List(); len(records)>0{

		for _, jobRecord := range records{
			job := Job.New("")
			if e := job.UnSerialize(jobRecord); e != nil{
				slog.Debugf("[storage.Service -> List] Skipped: ", e)
				continue
			}
			jobList = append(jobList, job)
		}
	}

	return jobList
}


func (s service)Lock(job api.Job) bool{
	return s.adapter.Lock(job.Index())
}

func (s service)UnLock(job api.Job) bool{
	return s.adapter.UnLock(job.Index())
}

func (s service)Flush(){
	s.adapter.Flush()
}

func (s service)Version() string{
	return s.adapter.Version()
}

func (s service)Close(){
	s.adapter.Disconnect()
}