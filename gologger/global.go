package gologger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const LEVEL_TRACE int = 0 // Tracking execution flow, detailed internal state.
const LEVEL_DEBUG int = 1 // Development and troubleshooting.
const LEVEL_INFO int = 2  // Application startup, user actions, expected events.
const LEVEL_WARN int = 3  // Deprecated API usage, high memory usage.
const LEVEL_ERROR int = 4 // Failed database query, network timeouts
const LEVEL_FATAL int = 5 // Data corruption, unrecoverable crash.

const MAXIMUM_LOG_FILE_BYTE_COUNT int = 10 * 1024 * 1024 // 10 MB

var isPrinting bool
var isLogging bool
var logFileDirectory string
var logFileObject *os.File
var logFileMutex sync.Mutex
var logFileByteCount int

func rotateLogFiles() {
	timestamp := time.Now().Format(time.RFC3339)
	_, err := os.Stat(logFileDirectory)
	if os.IsNotExist(err) {
		if isPrinting {
			noLogDirectoryMessage := fmt.Sprintf("(%s) [INFO] Log's directory does not exist. Creating ( %s ).", timestamp, logFileDirectory)
			fmt.Println(noLogDirectoryMessage)
		}
		err := os.Mkdir(logFileDirectory, 0755)
		if err != nil {
			if isPrinting {
				logDirectoryMkdirFailMessage := fmt.Sprintf("(%s) [WARN] Could not create the log's directory at ( %s ). The application will not write logs to the log file.", timestamp, logFileDirectory)
				fmt.Println(logDirectoryMkdirFailMessage)
			}
			isLogging = false
			logFileObject = nil
			return
		}
	} else if err != nil {
		if isPrinting {
			noLogDirectory := fmt.Sprintf("(%s) [WARN] Could not check for the log's directory existance at ( %s ). The application will not write logs to the log file.", timestamp, logFileDirectory)
			fmt.Println(noLogDirectory)
		}
		return
	}
	fileID := generateFileID()
	loggingFile := filepath.Join(logFileDirectory, fmt.Sprintf("%s.jsonl", fileID))
	logFileMutex.Lock()
	defer logFileMutex.Unlock()
	logFileObject, err = os.OpenFile(loggingFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if isPrinting {
			logFileFail := fmt.Sprintf("(%s) [WARN] Could not open the logs file at ( %s ). The application will not write logs to the log file.", timestamp, loggingFile)
			fmt.Println(logFileFail)
		}
		isLogging = false
		logFileObject = nil
		return
	}
	isLogging = true
}

func Initialize(printing bool, logging bool, logsDirectory string) {
	isPrinting = printing
	if !logging {
		return
	}
	timestamp := time.Now().Format(time.RFC3339)
	absoluteLogsDirectory, err := filepath.Abs(logsDirectory)
	if err != nil {
		logDirectoryAbsFailMessage := fmt.Sprintf("(%s) [ERROR] Could not get the absolute path of log's directory from ( %s ). The application will not write logs to the log's file.", timestamp, logsDirectory)
		fmt.Println(logDirectoryAbsFailMessage)
		return
	}
	logFileDirectory = absoluteLogsDirectory
	rotateLogFiles()
}

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

func LogWithExtra(level int, message string, extraFields map[string]string) {
	if !isLogging && !isPrinting {
		return
	}

	levelTag := getLevelTag(level)
	timestamp := time.Now().Format(time.RFC3339)

	if isPrinting && message != "" {
		logString := fmt.Sprintf("(%s) [%s] %s", timestamp, levelTag, message)
		fmt.Println(logString)
	}

	if isLogging && logFileObject != nil {
		logMap := map[string]string{
			"level":     levelTag,
			"timestamp": timestamp,
			"message":   message,
		}

		for key, value := range extraFields {
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

		logFileMutex.Lock()
		defer logFileMutex.Unlock()
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
