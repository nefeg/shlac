package cli

import sapi "shlacd/app/api"


type Context interface{
	List() []sapi.Job
	Get(job sapi.Job) sapi.Job

	Add(job sapi.Job, force bool) bool
	Import(cronLine string, checkDuplicates bool) (bool)


	Remove(job sapi.Job) bool

	Purge()
	Term()
}
