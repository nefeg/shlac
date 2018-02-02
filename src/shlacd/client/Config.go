package client

import (
	"log"
	"shlacd/client/socket"
	"shlacd/cli"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

type Config struct {

	Type    string `json:"type"`

	Options struct {
		Network string `json:"network"`
		Address string `json:"address"`
		Path    string `json:"path"`
	} `json:"options"`
}

func Resolve(conf Config) (client Handler){

	switch conf.Type {
		case "socket":
			slog.DebugLn("[client.config] Resolve: socket")
			client = Handler( socket.NewHandler( socket.NewConnectionConf(conf.Options.Network, conf.Options.Address),  cli.New() ))

		default:
			log.Fatalln("[client.config]Resolve(panic): Unknown client type: ", conf.Type)
	}

	return client
}