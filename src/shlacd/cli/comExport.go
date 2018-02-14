package cli

import (
	"github.com/urfave/cli"
	"errors"
	"fmt"
)

func NewComExport(context *Context) cli.Command {

	return cli.Command{
		Name:    "export",
		Aliases: []string{"x"},
		Usage:   "Export list of jobs",
		UsageText: "" +
			"\tlist [--full] [--no-legacy]",


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

			options := viewOptions{c.Bool("full"), !c.Bool("no-legacy")}

			c.App.Writer.Write( []byte( view( (*context).List(), options ) ) )

			return err
		},
	}
}
