package logger

import (
	"log"
	"os"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger interface {
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
}

var logLevel = DEBUG

// suggest use in system init
func SetLogLevel(level LogLevel) {
	logLevel = level
}

func GetLogger(module string) Logger {
	return NewGoLogger(module, logLevel)
}

type goLogger struct {
	level  LogLevel
	logger *log.Logger
}

func NewGoLogger(prefix string, level LogLevel) *goLogger {
	logger := log.New(os.Stdout, "["+prefix+"]", log.LstdFlags)
	return &goLogger{
		level:  level,
		logger: logger,
	}
}

func (g *goLogger) Debugf(format string, v ...interface{}) {
	if g.level > DEBUG {
		return
	}
	g.logger.Printf(format, v...)
}

func (g *goLogger) Debug(v ...interface{}) {
	if g.level > DEBUG {
		return
	}
	g.logger.Print(v...)
}

func (g *goLogger) Infof(format string, v ...interface{}) {
	if g.level > INFO {
		return
	}
	g.logger.Printf(format, v...)
}

func (g *goLogger) Info(v ...interface{}) {
	if g.level > INFO {
		return
	}
	g.logger.Print(v...)
}

func (g *goLogger) Warnf(format string, v ...interface{}) {
	if g.level > WARN {
		return
	}
	g.logger.Printf(format, v...)
}

func (g *goLogger) Warn(v ...interface{}) {
	if g.level > WARN {
		return
	}
	g.logger.Print(v...)
}

func (g *goLogger) Errorf(format string, v ...interface{}) {
	if g.level > ERROR {
		return
	}
	g.logger.Printf(format, v...)
}

func (g *goLogger) Error(v ...interface{}) {
	if g.level > ERROR {
		return
	}
	g.logger.Print(v...)
}
