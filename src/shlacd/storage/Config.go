package storage

import (
	"shlacd/app/api"
	"shlacd/storage/adapters"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

// storage config
type Config struct {
	Type    string `json:"type"`

	Options struct {
		Network string `json:"network"`
		Address string `json:"address"`
		Key     string `json:"key"`
		Path    string `json:"path"`
	} `json:"options"`
}

func Resolve(conf Config) (storage api.Storage){

	var adapter Adapter

	slog.Debugln("[storage->Resolve] Resolving storage config...")


	switch conf.Type {
	case "redis":
		slog.Debugln("[storage->Resolve] Resolved: `redis`")
		adapter = adapters.NewRedisAdapter(conf.Options.Network, conf.Options.Address, conf.Options.Key)

	case "file":
		// todo implement this

	case "script":
		// todo implement this
	}

	if adapter == nil{
		slog.Panicln("Unknown storage adapter")
	}


	return New(adapter)
}
