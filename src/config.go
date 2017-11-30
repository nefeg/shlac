package main

import (
	"storage"
	"client"
	"executor"
)

type Config struct {

	// storage config
	Storage  storage.Config `json:"storage"`

	// client config
	Client client.Config `json:"client"`

	// executor config
	Executor executor.Config `json:"executor"`
}