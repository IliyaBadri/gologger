package gologger

import (
	"os"
	"sync"
)

const LEVEL_TRACE int = 0 // Tracking execution flow, detailed internal state.
const LEVEL_DEBUG int = 1 // Development and troubleshooting.
const LEVEL_INFO int = 2  // Application startup, user actions, expected events.
const LEVEL_WARN int = 3  // Deprecated API usage, high memory usage.
const LEVEL_ERROR int = 4 // Failed database query, network timeouts
const LEVEL_FATAL int = 5 // Data corruption, unrecoverable crash.

const MAXIMUM_LOG_FILE_BYTE_COUNT int = 10 * 1024 * 1024 // 10 MB
const LOG_QUEUE_SIZE = 100                               // Default

type standardLog struct {
	level       int
	message     string
	extraFields map[string]string
}

var logChannel chan *standardLog = make(chan *standardLog, LOG_QUEUE_SIZE)

var isPrinting bool = true // if true the logs should be printed to the console.
var isLogging bool = true  // if true the logs should be written to the log file.

var logFileDirectory string = ""
var logFileObject *os.File = nil
var logFileByteCount int = 0

var noNewLogs bool
var noNewLogsMutex sync.Mutex

var fullStopMutex sync.Mutex
var fullStopSignal = sync.NewCond(&fullStopMutex)
