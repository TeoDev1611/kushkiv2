package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	DebugMode bool
	logFile   *os.File
)

// Init initializes the file logger.
func Init() error {
	// Obtener directorio de trabajo o configuración
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	
	logPath := filepath.Join(cwd, "kushki_app.log")
	
	// Abrir archivo en modo Append
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	
	logFile = f
	return nil
}

// Close closes the log file.
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func writeToLog(level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, msg)

	// 1. Escribir en Consola (Stdout)
	// Para consola usamos un formato más corto
	consolePrefix := fmt.Sprintf("[%s %s] ", level, time.Now().Format("15:04:05"))
	fmt.Print(consolePrefix + msg + "\n")

	// 2. Escribir en Archivo
	if logFile != nil {
		if _, err := logFile.WriteString(logLine); err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
		}
	}
}

// Debug prints messages only if DebugMode is true.
func Debug(format string, args ...interface{}) {
	if DebugMode {
		writeToLog("DEBUG", format, args...)
	}
}

// Info prints standard informational messages.
func Info(format string, args ...interface{}) {
	writeToLog("INFO", format, args...)
}

// Error prints error messages.
func Error(format string, args ...interface{}) {
	writeToLog("ERROR", format, args...)
}

