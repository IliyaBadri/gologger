# gologger

**gologger** is a simple and flexible logging package for Go, designed for both development and production environments. It supports logging at different levels, printing to the console, writing logs to a file, and rotating log files once a certain size limit is reached.

**WARNING:** This code is not bulletproof and is not field-tested.

## Features

-   **Log Levels**: Logs can be categorized into different levels:
    
    -   `TRACE`: For tracking execution flow and detailed internal state.
    -   `DEBUG`: For development and troubleshooting.
    -   `INFO`: For application startup, user actions, and expected events.
    -   `WARN`: For deprecated API usage and warnings like high memory usage.
    -   `ERROR`: For failed operations such as database queries and network timeouts.
    -   `FATAL`: For critical issues like data corruption and unrecoverable crashes.
-   **File-based Logging**: Logs can be written to a file, with automatic log rotation when the file exceeds a size threshold (default 10 MB).
    
-   **Thread-safe**: The logger operates safely across multiple goroutines, thanks to its internal synchronization mechanisms.
    
-   **Customizable Output**: You can enable or disable printing to the console or logging to a file.

- **It's all self contained**: No external dependencies involved.
    

## Usage

### Initialization

Before using GoLogger, you need to put the gologger package from the repository in your project's directory.

Here's a basic Example on how to use it:
```go
package main

import (
	"<yourproject>/gologger"
)

func main() {
	// Initialize with:
	// printing enabled, logging enabled, and the log file directory
	gologger.Initialize(true, true, "./logs")

	// You can start logging messages right after initialization.
	gologger.Log(gologger.LEVEL_INFO, "Application started")
	
	// Never forget to stop the logger before exiting
	gologger.WaitStopLogs()
}
```

### Logging a Message

Once initialized, you can log messages at different levels using the `Log` or `LogWithExtraFields` functions.

#### Basic Logging

```go
gologger.Log(gologger.LEVEL_ERROR, "An error occurred")
```

#### Logging with Extra Fields

```go
fields := map[string]string{
	"userID": "12345",
	"action": "login",
}
gologger.LogWithExtraFields(gologger.LEVEL_INFO, "User action logged", fields)
```

- For now the logs with extra fields are only written to the log file.

### Stopping Logging

If you need to stop logging (e.g., on program shutdown), you can use the following function to stop receiving new logs and wait for the existing logs to be written:

```go
gologger.WaitStopLogs()
```
## Other features
### File Rotation

GoLogger automatically rotates log files when the log file size exceeds the defined limit (10 MB by default). The log files are saved in the specified directory in JSON Lines format.

### Log File Directory

By default, GoLogger looks for the log directory. If the directory doesn't exist, GoLogger will attempt to create it. You can customize this behavior by passing a different path to the `Initialize` function.

### Thread Safety

The GoLogger package ensures thread safety for all logging operations, so you can safely log messages from multiple goroutines without worrying about race conditions.

### File Format

The log entries are saved in JSON Lines format, which makes it easy to parse and process the logs. Each log entry looks like this:

```json
{
  "level": "INFO",
  "timestamp": "2025-03-06T10:20:30Z",
  "message": "Application started",
  "userID": "12345",
  "action": "login"
}

```

### Configuration

-   **Console Output**: You can control whether logs are printed to the console by setting the `isPrinting` flag during initialization.
-   **File Output**: You can control whether logs are written to a file by setting the `isLogging` flag.
-   **Log File Size Limit**: When the log file exceeds 10 MB, the logger automatically rotates the log file.

### Error Handling

If any errors occur during file handling (such as file creation or writing errors), they will be logged to the console if `isPrinting` is enabled. Additionally, log writing is halted if the log file can't be accessed.

## License

GoLogger is released under the MIT License. See LICENSE for details.

----------

Feel free to open issues or contribute to the repository!
