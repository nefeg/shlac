package api


type Storage interface {

	Connect() (isConnected bool)
	Disconnect()

	Exists(index string) bool
	Add(index string, record string, force bool) (result bool, err error)
	Get(index string) (record string)
	Rm(index string) bool
	List() (data map[string]string)

	Lock(index string) bool
	UnLock(index string)

	Version() (version string)
	Flush()
}