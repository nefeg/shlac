package client

import (
	"log"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"shlacd/client/socket"
	"shared/config/addr"
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
			connection := &addr.Config{
				Protocol:conf.Options.Network,
				Address:conf.Options.Address,
			}

			slog.DebugF("[client.config] Resolve: socket [%s]\n", connection)

			client = Handler( socket.NewHandler( connection ))

		default:
			log.Fatalln("[client.config]Resolve(panic): Unknown client type: ", conf.Type)
	}

	return client
}