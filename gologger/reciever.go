package gologger

func LogWithExtraFields(level int, message string, extraFields map[string]string) {
	if !isNewLogsAllowed() {
		return
	}
	logObject := standardLog{
		level:       level,
		message:     message,
		extraFields: extraFields,
	}
	logChannel <- &logObject
}

func Log(level int, message string) {
	if !isNewLogsAllowed() {
		return
	}
	emptyFields := map[string]string{}
	logObject := standardLog{
		level:       level,
		message:     message,
		extraFields: emptyFields,
	}
	logChannel <- &logObject
}
