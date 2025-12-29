package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	file *os.File
}

func New() (*Logger, error) {
	logDir := "output"
	logFile := filepath.Join(logDir, fmt.Sprintf("scan_report_%s.log", time.Now().Format("20060102_150405")))

	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &Logger{file: f}, nil
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[INFO] "+format, v...)
	log.Println(msg)
	if l.file != nil {
		fmt.Fprintf(l.file, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[ERR] "+format, v...)
	log.Println(msg)
	if l.file != nil {
		fmt.Fprintf(l.file, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	}
}

func (l *Logger) Success(format string, v ...interface{}) {
	msg := fmt.Sprintf("[SUCCESS] "+format, v...)
	log.Println(msg)
	if l.file != nil {
		fmt.Fprintf(l.file, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("[WARN] "+format, v...)
	log.Println(msg)
	if l.file != nil {
		fmt.Fprintf(l.file, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	}
}
