package gologger

func isNewLogsAllowed() bool {
	noNewLogsMutex.Lock()
	newLogsAllowed := !noNewLogs
	noNewLogsMutex.Unlock()
	return newLogsAllowed
}

func WaitStopLogs() {
	noNewLogsMutex.Lock()
	noNewLogs = true
	noNewLogsMutex.Unlock()
	fullStopMutex.Lock()
	fullStopSignal.Wait()
	fullStopMutex.Unlock()
}
