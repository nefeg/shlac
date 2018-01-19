package api


type Executor interface {

	Exec(jobs ...Job) (outs [][]byte, errs []error)

	OnStart(f ...func(jobs []Job, err error))
	OnComplete(f ...func(jobs []Job, err error))

	OnItemStart(f ...func(job Job, err error, out []byte) bool )
	OnItemComplete(f ...func(job Job, err error, out []byte))
}