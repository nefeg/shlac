package api

import "time"

type Job interface {
	Id()                string
	Command()           string
	CronLine()          string
	TimeStart(fromTime time.Time) time.Time

	SetCommand(command string)
	SetCronLine(timeLine string)

	Serialize() string
	UnSerialize(data string)
}
