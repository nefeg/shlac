package cli

import (
	"github.com/urfave/cli"
	"errors"
	"fmt"
)

func NewComPurge(context *Context) cli.Command {

	return cli.Command{
		Name:    "purge",
		Usage:   "Remove all job",
		UsageText: "" +
			"\tpurge",

		Action:  func(c *cli.Context) (err error) {

			defer func(err *error){
				if r := recover(); r != nil{
					*err = errors.New(fmt.Sprintf("%s", r))
				}
			}(&err)

			(*context).Purge()

			return err
		},
	}
}

