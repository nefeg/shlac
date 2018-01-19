package Com

import (
	"shlacd/app/api"
	"errors"
)

type Quit struct{
	Com
}

const usageQuit = "usage: \n\t  quit (\\q) \n"

var ErrConnectionClosed = errors.New("** command <QUIT> received")

func (c *Quit)Exec(Tab api.TimeTable, args []string)  (string, error){

	return "OK", ErrConnectionClosed
}

func (c *Quit) Usage() string{
	return c.Desc() + "\n\t" + usageQuit
}