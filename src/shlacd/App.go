package main

type Application interface{
	Run()
	Stop(code int, message interface{})
}
