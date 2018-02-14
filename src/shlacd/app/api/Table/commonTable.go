package Table

import (
	"shlacd/app/api"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"errors"
	"fmt"
)

type commonTable struct {

	db          api.Storage
	jobs        []api.Job
	version     string
}


// constructor
func New (s api.Storage) *commonTable{

	t := &commonTable{ db:s }
	t.load()

	return t
}

func (t *commonTable) FindJob(needle api.Job) (api.Job){

	t.sync()
	for _, job := range t.jobs{
		if job.Index() == needle.Index() { return job }
	}

	return nil
}

func (t *commonTable) RmJob(job api.Job) bool{

	r := t.db.Rm(job)
	t.sync()

	return r
}

func (t *commonTable) AddJob(job api.Job, force bool) bool{

	defer func(job api.Job){

		if r := recover(); r!=nil{
			slog.Infoln("[Table -> AddJob] (panic): ", r)
			panic(r)
		}

		slog.Infof("[Table->AddJob] Successful added: Job#%s\n", job.Index())

		t.sync()

	}(job)


	if t.FindJob(job) != nil && !force{
		panic("Job already exist")
	}


	_, err := t.db.Add(job, force)
	if  err != nil{
		panic(err)
	}

	return true // todo fix it

	//t.jobs = append(t.jobs, Job)

}

func (t *commonTable) PullJob(job api.Job) (api.Job, error){

	var err error

	slog.Debugln("[Table -> PullJob] PullJob: Trying to lock Job...")
	if t.db.Lock(job){

		slog.Debugf("[Table -> PullJob] PullJob: Job #%s locked\n", job.Index())
		job = t.FindJob(job)
		if job == nil {
			t.db.UnLock(job)
			err = errors.New(fmt.Sprintf("Job '%s' not found", job.Index()))
			slog.Debugln("[Table -> PullJob] PullJob: ", err)
		}

	}else{
		err = errors.New(fmt.Sprintf("Locking for Job#%s fail", job.Index()))
		slog.Debugln("[Table -> PullJob] PullJob: ", err)
	}

	return job, err
}

func (t *commonTable) PushJob(job api.Job)  {

	slog.Debugf("[Table -> PushJob] Release lock for Job#%s\n", job.Index())

	r := t.db.UnLock(job)

	slog.Debugf("[Table -> PushJob] Release lock for Job#%s: %v \n", job.Index(), r)

}

func (t *commonTable) ListJobs() []api.Job{

	t.sync()

	return t.jobs
}

func (t *commonTable) Flush() {

	t.db.Flush()
	t.sync()
}

func (t *commonTable) Close(){
	t.db.Close()
}



func (t *commonTable) sync(){

	if !t.isSynced(){
		slog.Infof("[Table->sync] update version: %v --> %v\n", t.version, t.db.Version())
		t.load()
	}
}

func (t *commonTable) isSynced() bool{
	return t.version == t.db.Version()
}

func (t *commonTable) load(){

	t.jobs      = nil
	t.version   = t.db.Version()
	t.jobs      = t.db.List()

	// todo move to storage
	//for _, jobData := range t.db.List() {
	//	job := Job.New("")
	//	if e := job.UnSerialize(string(jobData)); e != nil{
	//		slog.Errln("[Table->load] ERR: Job skipped with error: ", e.Error())
	//		continue
	//	}
	//
	//	t.jobs = append(t.jobs, job)
	//}
}