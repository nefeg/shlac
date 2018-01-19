package cli

import "shlacd/app/api"

type Cmd interface {

	Resolve(cmdName string) (cmd string, args []string, err error)
	Exec(Tab api.TimeTable, args []string)  (string, error)
	Usage() string
}