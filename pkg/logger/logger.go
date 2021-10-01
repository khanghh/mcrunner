package logger

import (
	"fmt"
	"log"
)

type LogType int

var (
	LogLevel LogType = 0
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[0:37m"
	White  = "\033[37m"
)

const (
	debugPrefix   = "[DEBUG]"
	errorPrefix   = "[ERROR]"
	warningPrefix = "[WARN]"
	infoPrefix    = "[INFO]"
	spaceElem     = " "
)

const (
	Error LogType = iota // 0
	Warning
	Debug
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func Printf(format string, args ...interface{}) {
	log.SetPrefix(White)
	log.Print(infoPrefix+spaceElem, fmt.Sprintf(format, args...))
}

func Println(msg string) {
	log.SetPrefix(White)
	log.Println(infoPrefix, msg)
}

func Debugf(format string, args ...interface{}) {
	if LogLevel >= Debug {
		log.SetPrefix(Gray)
		log.Print(infoPrefix+spaceElem, fmt.Sprintf(format, args...))
	}
}

func Debugln(msg string) {
	if LogLevel >= Debug {
		log.SetPrefix(Gray)
		log.Println(debugPrefix, msg)
	}
}

func Warnf(format string, args ...interface{}) {
	if LogLevel >= Warning {
		log.SetPrefix(Yellow)
		log.Print(warningPrefix+spaceElem, fmt.Sprintf(format, args...))
	}
}

func Warnln(msg string) {
	if LogLevel >= Warning {
		log.SetPrefix(Yellow)
		log.Println(warningPrefix, msg)
	}
}

func Errorln(v ...interface{}) {
	if LogLevel >= Error {
		log.SetPrefix(Red)
		v = append([]interface{}{errorPrefix}, v...)
		log.Println(v...)
	}
}

func Fatalln(v ...interface{}) {
	log.SetPrefix(Red)
	v = append([]interface{}{errorPrefix}, v...)
	log.Fatalln(v...)
}
