package api

type TimeTable interface {

	FindJob(jobId string) (job Job)

	AddJob(job Job, force bool)

	RmJob(jobId string) bool

	PullJob(jobId string) (job Job)

	PushJob(job Job)

	ListJobs() []Job

	Flush()
	Close()
}




