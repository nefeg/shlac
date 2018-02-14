package Context

import (
	sapi "shlacd/app/api"
	"strings"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"regexp"
	"shlacd/app/api/Job"
)

func New(T sapi.Table) *context{


	return &context{T}
}

type context struct {

	table   sapi.Table
}


func (c *context) Import(line string, checkDuplicates bool) (result bool) {

	slog.Debugf("[cli.Context -> Import] Data: `%s`\n", line)
	slog.Debugf("[cli.Context -> Import] checkDuplicates: `%v`\n", checkDuplicates)

	if line != "" && line != "\n"{

		re := regexp.MustCompile(`^([0-9,\/*-LW#]+\s+){5,7}(.+)$`)
		matches := re.FindStringSubmatchIndex(line)

		slog.Debugln("[cli.Context -> Import] matches: ", matches)

		if len(matches) < 4{
			slog.Debugln("[cli.Context -> Import] parse: incorrect format")
			panic("incorrect format of imported data")
		}


		cronLine    := strings.Trim( line[:matches[3]], " \t" )
		commandLine := strings.Trim( line[matches[3]:], " \t" )

		slog.Debugf("[cli.Context -> Import] parsed `cronLine`: `%s`\n", cronLine)
		slog.Debugf("[cli.Context -> Import] parsed `commandLine`: `%s`\n", commandLine)


		importJob := Job.New("")
		importJob.CommandX(commandLine)
		//importJob.CommandX(fmt.Sprintf(`%q`, commandLine))
		importJob.TimeLineX(cronLine)

		slog.Debugln("[cli.Context -> Import] Job: ", importJob)

		if cronLine[:1] == `#`{
			slog.Infof("[cli.Context -> Import] SKIPP (disabled)>> %s\n", importJob)
			return
		}

		if checkDuplicates && c.isDuplicated(importJob.String()){
			slog.Infof("[cli.Context -> Import] SKIPP (duplicated)>> %s\n", importJob)
			return
		}

		slog.Infof("[cli.Context -> Import] IMPORT>> %s\n", importJob)

		result = c.Add(importJob, true)
	}

	slog.Debugf("[cli.Context -> Import] Result: %v\n", result)


	return result
}


func (c *context) List() []sapi.Job {

	return c.table.ListJobs()
}


func (c *context) Get(job sapi.Job) sapi.Job{

	return c.table.FindJob(job)
}

func (c *context) Add(job sapi.Job, force bool) bool{

	return c.table.AddJob(job, force)
}

func (c *context) Remove(job sapi.Job) bool{

	return c.table.RmJob(job)
}

func (c *context) Purge(){
	c.table.Flush()
}

func (c *context) Term(){
	c.table.Close()
}


func (c *context) isDuplicated(needle string) (r bool) {

	slog.Debugf("[cli.Context -> isDuplicated] In(cronLine): `%s`\n", needle)

	ws := regexp.MustCompile(`\s+`)
	needle = ws.ReplaceAllString(needle, "")

	slog.Debugf("[cli.Context -> isDuplicated] Normalized: `%s`\n", needle)


	current := ""
	for _,job := range c.List(){

		current = ws.ReplaceAllString(job.String(),"")

		slog.Debugf("[cli.Context -> isDuplicated] Compare (needle, current): \n\t-`%s`\n\t-`%s`\n", needle, current)
		if current == needle{
			r = true
		}
	}

	slog.Debugf("[cli.Context -> isDuplicated] Result: `%v`\n", r)

	return r
}
