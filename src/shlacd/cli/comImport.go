package cli

import (
	"github.com/urfave/cli"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"strings"
	"bufio"
	"fmt"
	"io"
	"errors"
	"os"
)

func NewComImport(context *Context) cli.Command {

	return cli.Command {//IMPORT
		Name:    "import",
		Aliases: []string{"i"},
		Usage:   "import jobs from cron-formatted file",
		UsageText: "" +
			"\tshlac import [--purge] [--skip-check] \"<cron-formatted line>\"\n" +
			"\tshlac import [--purge] [--skip-check] -f <path/to/import/file>\n" +
			"\tcat <path/to/import/file> | shlac import  [--purge] [--skip-check]" +
			"",

		Flags: 	[]cli.Flag{
			cli.BoolFlag{
				Name:  "purge",
				Usage: "delete jobs before import",
			},

			cli.BoolFlag{
				Name:  "skip-check, s",
				Usage: "add Job even if same is already exist (skip checking for duplicates)",
			},

			cli.StringFlag{
				Name:  "file, f",
				Usage: "import from file",
			},
		},

		Action:  func(c *cli.Context) (err error) {

			var records []string

			// clean Table before import
			if c.Bool("purge") {
				(*context).Purge()
			}


			//############## File ######################
			// Import from file
			if path := c.String("file"); path != "" {

				slog.DebugF("[cli.import] Source: file\n")

				if source, err := os.Open(path); err == nil {
					// !get cron-lines
					records = scan(source)
					source.Close()
				}else{
					slog.DebugF("[cli.import] Source error: %s", err)
					return err
				}

			//##########################################



			//############## COMMAND ARGS ##############
			}else {
				// check COMMAND ARGS
				cronLine := strings.Join(c.Args(), " ")
				slog.DebugF("[cli.import] cron-line: `%s`\n", cronLine)

				if cronLine != "" { // import from args
					records = append(records, cronLine)
				}
			}
			//##########################################



			//############## STDIN #####################
			// or import from stdIn/pipe
			if len(records) == 0 {
				c.App.Writer.Write([]byte("Interactive mode\nPress ^D(Ctrl-D) for end or ^C(Ctrl-C) for halt program\n"))

				slog.DebugF("[cli.import] Source: stdin\n")

				// !get cron-lines
				var source io.Reader = os.Stdin
				records = scan(source)
			}
			//##########################################



			//############# IMPORT #####################
			var imported int
			if len(records) > 0 {
				for _,line := range records{
					if (*context).Import(line, !c.Bool("skip-check")){
						c.App.Writer.Write([]byte(fmt.Sprintf("***IMPORT\t-->`%s`\n", line)))
						imported ++
					}else{
						c.App.Writer.Write([]byte(fmt.Sprintf("***SKIP\t-->`%s`\n", line)))

					}
				}

			// if no records imported from file/stdin and no cron-lines found in args then start to panic
			}else{
				err = errors.New("ERR: expected data for import")
			}
			//##########################################


			// show stats
			c.App.Writer.Write( []byte(fmt.Sprintf("===Total: %d record(s) imported\n", imported)) )

			return err
		},
	}
}

func scan(source io.Reader) (records []string){

	scanner     := bufio.NewScanner(source)

	for scanner.Scan() {
		records = append(records, scanner.Text())
	}

	return records
}

