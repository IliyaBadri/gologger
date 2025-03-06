package gologger

import (
	"encoding/json"
	"fmt"
	"time"
)

func getLevelTag(level int) string {
	clamped_level := level
	if clamped_level < 0 {
		clamped_level = 0
	}
	if clamped_level > 5 {
		clamped_level = 5
	}
	switch clamped_level {
	case LEVEL_TRACE:
		return "TRACE"
	case LEVEL_DEBUG:
		return "DEBUG"
	case LEVEL_INFO:
		return "INFO"
	case LEVEL_WARN:
		return "WARN"
	case LEVEL_ERROR:
		return "ERROR"
	case LEVEL_FATAL:
		return "FATAL"
	}
	return "INFO"
}

func logWriterRoutine() {
	for {
		if !isNewLogsAllowed() && len(logChannel) <= 0 {
			fullStopMutex.Lock()
			fullStopSignal.Broadcast()
			fullStopMutex.Unlock()
			break
		}
		lastLog := <-logChannel
		if !isLogging && !isPrinting {
			continue
		}
		levelTag := getLevelTag(lastLog.level)
		timestamp := time.Now().Format(time.RFC3339)
		if isPrinting && lastLog.message != "" {
			logString := fmt.Sprintf("(%s) [%s] %s", timestamp, levelTag, lastLog.message)
			fmt.Println(logString)
		}
		if isLogging && logFileObject != nil {
			logMap := map[string]string{
				"level":     levelTag,
				"timestamp": timestamp,
				"message":   lastLog.message,
			}

			for key, value := range lastLog.extraFields {
				if key == "level" || key == "timestamp" || key == "message" {
					continue
				}
				logMap[key] = value
			}

			jsonLog, err := json.Marshal(logMap)
			if err != nil {
				if isPrinting {
					logFileFail := fmt.Sprintf("(%s) [ERROR] An error occured while trying to marshal a log map into JSON format. This log will not be written to the log file.\n%s", timestamp, err)
					fmt.Println(logFileFail)
				}
				return
			}

			bytesToWrite := append(jsonLog, []byte("\n")...)
			bytesWritten, err := logFileObject.Write(bytesToWrite)
			if err != nil {
				if isPrinting {
					logFileFail := fmt.Sprintf("(%s) [ERROR] An error occured while trying to write a log to the log file. This log will not be written to the log file.\n%s", timestamp, err)
					fmt.Println(logFileFail)
				}
			}

			logFileByteCount += bytesWritten

			if logFileByteCount >= MAXIMUM_LOG_FILE_BYTE_COUNT {
				rotateLogFiles()
				logFileByteCount = 0
			}
		}
	}
}
