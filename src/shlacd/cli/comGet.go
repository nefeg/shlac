package cli

import (
	"github.com/urfave/cli"
	"errors"
	"fmt"
	"shlacd/app/api/Job"
)

// TODO merge `get` and `list` to `export`
func NewComGet(context *Context) cli.Command {

	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "Get job by id",
		UsageText: "" +
			"\tget <id>",

		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "full,f",
				Usage: "Show full information about the job",
			},
			cli.BoolFlag{
				Name:  "no-legacy,n",
				Usage: "Non cron-formatted response",
			},

		},

		Action:  func(c *cli.Context) (err error) {

			defer func(err *error){
				if r := recover(); r != nil{
					*err = errors.New(fmt.Sprintf("%s", r))
				}
			}(&err)

			jobIndex := c.Args().Get(0)
			if jobIndex == "" {
				panic("ERR: expected job index\nsee `get --help`")
			}


			options := viewOptions{c.Bool("full"), !c.Bool("no-legacy")}

			job := Job.New(jobIndex)
			c.App.Writer.Write( []byte( viewItem((*context).Get(job), options) ))

			return err
		},
	}
}
