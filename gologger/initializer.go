package gologger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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
	if logFileObject != nil {
		err = logFileObject.Close()
		if err != nil && isPrinting {
			logDirectoryAbsFailMessage := fmt.Sprintf("(%s) [ERROR] Could not close the previous log file before rotating to a new one.", timestamp)
			fmt.Println(logDirectoryAbsFailMessage)
		}
	}
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
	isLogging = logging
	if !isLogging {
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
	go logWriterRoutine()
}
