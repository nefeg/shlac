package Job

import (
	"time"
	"encoding/json"
	"log"
	"github.com/umbrella-evgeny-nefedkin/cronexpr"
)

type job struct {
	ID      string
	Cln     string // cron-line
	Cmd     string // command
}

func New() *job{
	return &job{}
}

func (j *job)Id() string{
	return j.ID
}

func (j *job)Command() string{
	return j.Cmd
}

func (j *job)CronLine() string{
	return j.Cln
}

func (j *job)TimeStart(fromTime time.Time) time.Time{

	return cronexpr.MustParse(j.CronLine()).Next(fromTime)
}


func (j *job)SetID(ID string){
	if j.Id() != ""{
		log.Panicf("[job] SetID(panic): try to change id %s --> %s \n", j.Id(), ID)
	}

	j.ID = ID
}

func (j *job)SetCommand(command string) {
	j.Cmd = command
}

func (j *job)SetCronLine(timeLine string){

	if _, e := cronexpr.Parse(timeLine); e != nil{
		log.Panicln("Invalid format for '-cron'")
	}

	j.Cln = timeLine
}



// interface storage.Serializable
func (j *job) Serialize() string{

	s, err := json.Marshal(j)
	if err != nil{
		log.Fatalln(err)
	}

	log.Println("[Job]Serialize: ", string(s), err)

	return string(s)
}

func (j *job) UnSerialize(data string){

	err := json.Unmarshal([]byte(data), j)
	if err != nil{
		log.Fatalln(err)
	}

	log.Println("[Job]Unserialize: ", j, err)
}

