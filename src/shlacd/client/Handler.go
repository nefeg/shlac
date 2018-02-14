package client

import (
	"shlacd/cli"
)

type Handler interface {

	Handle(ctx cli.Context)
}

