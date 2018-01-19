package app

import (
	"shlacd/storage"
	"shlacd/client"
	"shlacd/executor"
)

type Config struct {

	// storage config
	Storage  storage.Config `json:"storage"`

	// client config
	Client client.Config `json:"client"`

	// executor config
	Executor executor.Config `json:"executor"`
}
