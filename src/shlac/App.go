package main

type App interface{

	Export(options interface{}) (str string)
	Import(filePath string, checkDuplicates bool)
	ImportLine(cronString string, checkDuplicates bool)
	Remove(jobId string)
	Purge()
}