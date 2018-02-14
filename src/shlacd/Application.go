package main

type Application interface{
	IsDebug() bool
	Run()
	Stop(code int, message interface{})
}
