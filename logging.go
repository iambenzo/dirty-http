package dirtyhttp

import (
	"log"
	"os"
)

type messageType int

const (
	INFO messageType = 0 + iota
	WARNING
	ERROR
	FATAL
)

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

func writeLog(messagetype messageType, message string) {
	switch messagetype {
	case INFO:
		log.SetPrefix(Cyan + "[INFO] " + Reset)
		log.Println(message)
	case WARNING:
		log.SetPrefix(Yellow + "[WARNING] " + Reset)
		log.Println(message)
	case ERROR:
		log.SetPrefix(Red + "[ERROR] " + Reset)
		log.Println(message)
	case FATAL:
		log.SetPrefix(Red + "[FATAL] " + Reset)
		log.Fatalln(message)
	}
}

type logger struct {}

func (l *logger) Info(message string) {
    writeLog(INFO, message)
}

func (l *logger) Warning(message string) {
    writeLog(WARNING, message)
}

func (l *logger) Error(message string) {
    writeLog(ERROR, message)
}

func (l *logger) Fatal(message string) {
    writeLog(FATAL, message)
    os.Exit(1)
}

