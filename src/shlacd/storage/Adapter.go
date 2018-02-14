package storage

type Adapter interface {

	Connect() (isConnected bool)
	Disconnect()

	Exists(index string) bool
	// by default assignedIndex == index, or generate by storage if index == ""
	Add(index string, record string, force bool) (err error)
	Get(index string) (record string)
	Rm(index string) bool
	List() (data map[string]string)

	Lock(index string) bool
	UnLock(index string)

	Version() (version string)
	Flush()
}
