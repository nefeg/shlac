package api

type Table interface {

	FindJob(needle Job) Job

	AddJob(job Job, force bool) bool

	RmJob(job Job) bool

	PullJob(job Job) (Job, error)

	PushJob(job Job)

	ListJobs() []Job

	Flush()
	Close()
}




