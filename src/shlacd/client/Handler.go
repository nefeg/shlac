package client

import (
	"shlacd/app/api"
)

type Handler interface {

	Handle(Tab api.TimeTable)
}

