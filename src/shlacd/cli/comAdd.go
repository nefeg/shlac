package cli

import (
	"github.com/urfave/cli"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"shlacd/app/api/Job"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

func NewComAdd(context *Context) cli.Command {

	return cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add job",
		UsageText: "" +
			"\tshlac add [-i <index>] [--force] -t <cron-time>  -e <command to execute>",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "cmd,e",
				Usage: "Command to execute (required)",
			},
			cli.StringFlag{
				Name:  "index,i",
				Usage: "Set job index manually (job index is unique)",
			},
			cli.StringFlag{
				Name:  "timeline, t",
				Usage: "Cron-formatted time line (e.g. \"*/2 * * *\")",
			},

			cli.BoolFlag{
				Name:  "force",
				Usage: "Override duplicates",
			},
		},


		Action: func(c *cli.Context) (err error) {

			newJob := Job.New("")

			defer func(err *error){
				if r := recover(); r != nil{
					*err = errors.New(fmt.Sprintf("%s", r))
					slog.ErrLn("[cli.add] FAIL while adding job: ", newJob.Index())
				}
			}(&err)


			// Set 'command' (option required)
			if cmd := c.String("cmd");  cmd == ""{
				panic(errors.New(`ERR: "-cmd" required` + "\nsee `add --help`"))
			}else{
				newJob.CommandX(cmd)
			}


			// Set 'timeline' (option required)
			if timeline := c.String("timeline") ; timeline == ""{
				panic(errors.New(`ERR: "--timeline" required\nsee "add --help""`))
			}else{
				newJob.TimeLineX(timeline)
			}


			// or generate index if it's empty
			if index := c.String("index"); index != "" {
				newJob.IndexX(index)

			}else{
				if uid, err := uuid.NewV4(); err == nil{
					newJob.IndexX(uid.String())
				}else{
					panic(err)
				}
			}


			slog.DebugLn("[cli.add] Trying add to table: ", newJob)

			// Add job to table
			if (*context).Add(newJob, c.Bool("force")){
				c.App.Writer.Write( []byte(newJob.Index()) )
				slog.InfoLn("[cli.add] Added to table: ", newJob.Index())

			}else{
				panic("(*context).Add(..) return FALSE ")
			}

			return err
		},
	}
}
