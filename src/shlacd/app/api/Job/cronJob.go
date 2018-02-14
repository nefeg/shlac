package Job

import (
	"time"
	"encoding/json"
	"github.com/umbrella-evgeny-nefedkin/cronexpr"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"errors"
	"fmt"
)

type cronJob struct {
	index        string
	timeLine     string // cron-line
	command      string // command
}

type cronJobTemplate struct {
	Index       string      `json:"index"`
	TimeLine    string      `json:"timeLine"`
	Command     string      `json:"command"`
}

func New(index string) *cronJob{

	c := &cronJob{}
	c.IndexX(index)

	return c
}


// ID
func (c *cronJob)Index() string{
	return c.index
}

func (c *cronJob)IndexX(index string){
	//if c.index != "" || index == ""{
	//	err := fmt.Sprintf("try to change id `%s` --> `%s`", c.index, index)
	//	slog.Debugln("[Job->IndexX] (panic): ", err)
	//	panic(err)
	//}

	slog.Debugf("[Job->IndexX] try to change id `%s` --> `%s`\n", c.index, index)

	c.index = index
}


// command
func (c *cronJob)Command() string{
	return c.command
}

func (c *cronJob)CommandX(command string) {

	if command == "" || command == "\n"{
		err := "job command can't be empty"
		slog.Debugln("[Job->IndexX] (panic): ", err)
		panic(err)
	}

	c.command = command
}


func (c *cronJob)TimeLine() string{
	return c.timeLine
}


func (c *cronJob)TimeLineX(timeLine string){

	if _, e := cronexpr.Parse(timeLine); e != nil{
		err := fmt.Sprintf("Invalid timeline format: `%s`", timeLine)
		slog.Critf("[Job->TimeLineX] CRIT: %s\n", err)
		panic(err)
	}

	c.timeLine = timeLine
}


func (c *cronJob)TimeStart(fromTime time.Time) time.Time{

	return cronexpr.MustParse(c.TimeLine()).Next(fromTime)
}



// implementation of Marshaler/Unmarshaler interfaces
func (c *cronJob) MarshalJSON() ([]byte, error) {
	return json.Marshal(cronJobTemplate{
		Index:      c.index,
		Command:    c.command,
		TimeLine:   c.timeLine,
	})
}

func (c *cronJob) UnmarshalJSON(b []byte) (err error) {

	defer func(e *error){

		if r := recover(); r != nil{
			*e = errors.New(fmt.Sprint(r))
		}

	}(&err)


	temp := &cronJobTemplate{}

	err = json.Unmarshal(b, temp)
	if err == nil {
		c.IndexX(temp.Index)
		c.CommandX(temp.Command)
		c.TimeLineX(temp.TimeLine)
	}

	return err
}



// interface storage.Serializable
func (c *cronJob) Serialize() string{

	slog.Debugf("[Job->Serialize] Try to serialize Job: `%s`\n", c)

	s, err := json.Marshal(c)
	if err != nil{
		slog.Fatalln(err)
	}

	slog.Debugln("[Job->Serialize] Done: ", string(s), err)

	return string(s)
}

func (c *cronJob) UnSerialize(data string) (err error){

	slog.Debugf("[Job->Unserialize] Try to unserialize data: `%s`\n", data)

	err = json.Unmarshal([]byte(data), c)

	slog.Debugln("[Job->Unserialize] Done: ", c, err)

	return err
}


func (c *cronJob) String() string{

	cronLine := fmt.Sprintf("%s %s", c.TimeLine(), c.Command())

	//if c.Index() != "" {
	//	cronLine += " #" +c.Index()
	//}

	return cronLine
}

