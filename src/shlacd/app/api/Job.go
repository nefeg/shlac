package api

import "time"

type Job interface {
	Index()             string
	Command()           string
	TimeLine()          string

	IndexX(index string)
	CommandX(command string)
	TimeLineX(timeLine string)

	TimeStart(fromTime time.Time) time.Time

	Serialize() string
	UnSerialize(data string) (err error)

	String() string
}
