package api


type Storage interface {

	// by default assignedIndex == index, or generate by storage if index == ""
	Add(job Job, force bool) (Job, error)
	Load(job Job) (Job, error)
	Rm(job Job) bool
	List() []Job

	Lock(job Job) bool
	UnLock(job Job) bool

	Flush()

	Version() string

	Close()
}