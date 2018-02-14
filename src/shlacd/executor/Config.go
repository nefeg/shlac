package executor

import (
	"shlacd/app/api"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

// executor config
type Config struct {
	Type    string `json:"type"`

	Options struct {
		Silent  bool `json:"silent"`
		Async   bool `json:"async"`

	} `json:"options"`
}


func Resolve(conf Config) (exe api.Executor){

	switch conf.Type {
	case "local":
		exe = NewExecutorLocal(conf.Options.Silent, conf.Options.Async)

	default:
		slog.Panicln("Unknown client type")
	}

	return exe
}