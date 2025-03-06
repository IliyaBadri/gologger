package main

import "gologger/gologger"

func main() {
	gologger.Initialize(true, true, "./logs")
	gologger.Log(gologger.LEVEL_DEBUG, "This is a debugging message.")
	gologger.Log(gologger.LEVEL_INFO, "This is an informational log.")
	gologger.Log(gologger.LEVEL_FATAL, "This is a FATAL log!!")
	gologger.WaitStopLogs()
}
