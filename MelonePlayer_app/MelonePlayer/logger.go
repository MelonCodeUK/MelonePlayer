package MelonePlayer

// custom logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/inconshreveable/log15"
)

type Logger struct {
	logger log15.Logger
	file   *os.File
}

var nowt = time.Now()

// Create a new logger instance
func NewLogger(log_file_path string) Logger {
	logger := log15.New()

	logger.SetHandler(log15.StdoutHandler)
	file, err := os.OpenFile(log_file_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Crit("error opening log file: %v", err)
	}
	formattedTime := nowt.Format("[01-02|15:04:05]")
	file.WriteString(strings.Repeat("-", 12) + formattedTime + "-Writing to log file started!-" + strings.Repeat("-", 48) + "\n")
	fmt.Println(formattedTime + "-Writing to log file started!-" + strings.Repeat("-", 24))

	return Logger{
		logger: logger, file: file,
	}

}

// Info логирование
func (logger Logger) Write(str ...any) {
	message := fmt.Sprint(str...)
	fmt.Print(message)
	logger.writing("Write", message)

}

// Info логирование
func (logger Logger) Info(str ...any) {
	message := fmt.Sprint(str...)
	logger.writing("INFO", message)
	logger.logger.Info(message)
}

// Warn логирование
func (logger Logger) Warn(str ...any) {
	message := fmt.Sprint(str...)
	logger.writing("WARN", message)
	logger.logger.Warn(message)
}

// Debug логирование
func (logger Logger) Debug(str ...any) {
	if IsDebug {
		message := fmt.Sprint(str...)
		logger.writing("DEBUG", message)
		logger.logger.Debug(message)
	}
}

// Error логирование
func (logger Logger) Error(str ...any) {
	message := fmt.Sprint(str...)
	logger.writing("ERROR", message)
	logger.logger.Error(message)
}

// Fatal логирование
func (logger Logger) Fatal(str ...any) {
	message := fmt.Sprint(str...)
	logger.writing("FATAL", message)
	logger.logger.Crit(message)
}

// writing запись в файл
func (logger Logger) writing(type_error string, str string) {
	if IsSave_Log {
		formattedTime := nowt.Format("01-02|15:04:05")
		logEntry := fmt.Sprintf("%s %s %s", type_error, formattedTime, str)
		logger.file.WriteString(logEntry + "\n")
	}
}
