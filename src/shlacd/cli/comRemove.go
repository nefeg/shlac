package cli

import (
	"github.com/urfave/cli"
	"errors"
	"fmt"
	"shlacd/app/api/Job"
)

func NewComRemove(context *Context) cli.Command {

	return cli.Command{
		Name:    "remove",
		Aliases: []string{"r"},
		Usage:   "Remove job by index",
		UsageText: "" +
			"\tremove <index>",

		Action:  func(c *cli.Context) (err error) {

			defer func(err *error){
				if r := recover(); r != nil{
					*err = errors.New(fmt.Sprintf("%s", r))
				}
			}(&err)


			if c.Bool("all"){
				(*context).Purge()
				return err
			}


			job := Job.New(c.Args().Get(0))
			if job.Index() == "" {
				panic("ERR: expected job index\nsee `remove --help`")
			}

			(*context).Remove(job)

			return err
		},
	}
}
