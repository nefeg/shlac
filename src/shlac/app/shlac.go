package app

import (
	"fmt"

	"strings"
	"bytes"
	"regexp"
	"os"
	"bufio"
	"github.com/umbrella-evgeny-nefedkin/slog"
)



type shlac struct {
	sender  Sender
}


func New(sender Sender) *shlac{

	return &shlac{sender}
}


type ExportOpt struct {
	ShowId  bool
}

func (s *shlac) Export(options interface{}) (str string){

	response := string(s.sender.Send(`\l`)) // write command and read answer
	response = response[:len(response)-4] // remove terminal bytes

	if !options.(ExportOpt).ShowId{
		re := regexp.MustCompile(`(?m)^.+?\s+`) // remove jobs id
		str = re.ReplaceAllString(response, "")

	}else{

		matches := regexp.MustCompile(`(?m)^(.+?\s)([0-9,\/*-LW#]+\s+){5,7}(.+)$`).FindAllStringSubmatchIndex(response, -1)

		for i := range matches{
			str += strings.Trim(response[matches[i][3]:matches[i][7]] + ` # ` + response[matches[i][0]:matches[i][3]], " \t\n") + "\n"
		}
	}

	return str
}



func (s *shlac) Import(filePath string, checkDuplicates bool){

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner     := bufio.NewScanner(file)

	for scanner.Scan() {

		s.ImportLine(scanner.Text(), checkDuplicates)

		// MAIN LOOP
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}

func (s *shlac) ImportLine(cronString string, checkDuplicates bool){

	var importLine string

	if cronString == "\n"{
		importLine = cronString

	}else{
		re := regexp.MustCompile(`^([0-9,\/*-LW#]+\s+){5,7}(.+)$`)
		matches := re.FindStringSubmatchIndex(cronString)

		cronLine    := strings.Trim( cronString[:matches[4]], " \t" )
		commandLine := strings.Trim( cronString[matches[4]:], " \t" )

		importLine = fmt.Sprintf(`\a -cron "%s" -cmd %q`, cronLine, commandLine)
		slog.DebugLn("[app] ImportLine: ", importLine)

		if cronLine[:1] == `#`{
			fmt.Printf("SKIPP (disabled)>> %s\n", importLine)
			return
		}

		if checkDuplicates && s.isDuplicated(cronString){
			fmt.Printf("SKIPP (duplicated)>> %s\n", importLine)
			return
		}


		fmt.Printf("IMPORT>> %s\n", importLine)
	}


	s.sender.Send(importLine)
}


func (s *shlac) Purge(){
	s.sender.Send(`\r --all`)
}


func (s *shlac) Remove(jobId string){
	s.sender.Send(`\r -id `+jobId)
}

func (s *shlac) isDuplicated(cronLine string) bool {

	slog.DebugF("[app.shlac] isDuplicated (cronLine): `%s`\n", cronLine)

	cronLine = strings.Replace(cronLine, `"`, `\"`, -1)

	slog.DebugF("[app.shlac] isDuplicated (cronLine,normalized): `%s`\n", cronLine)

	response := s.sender.Send(`\g -c "` +cronLine+ `"`)

	slog.DebugLn("[app.shlac] isDuplicated (response,raw): ", response)
	slog.DebugF("[app.shlac] isDuplicated (response,string): `%s`\n", string(response))

	r := !bytes.Equal(response[:4], []byte{0, 10, 62, 62})

	slog.DebugLn("[app.shlac] isDuplicated (result): ", r)

	return r
}