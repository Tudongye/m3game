package log

import (
	"fmt"
	golog "log"
)

func newDefaultLogger() *defaultLogger {
	return &defaultLogger{
		lv: LogLvTrace,
	}
}

type defaultLogger struct {
	lv LogLv
}

func (l *defaultLogger) Output(depth int, logflag string, plus string, format string, v ...interface{}) {
	file, funcName, line := GetCallPath(depth)
	if plus == "" {
		golog.Printf("[%s] [%s:%d %s] %s\n", logflag, file, line, funcName, fmt.Sprintf(format, v...))
	} else {
		golog.Printf("[%s] [%s:%d %s][%s] %s\n", logflag, file, line, funcName, plus, fmt.Sprintf(format, v...))
	}
}

func (l *defaultLogger) SetLevel(level LogLv) {
	l.lv = level
}

func (l *defaultLogger) GetLevel() LogLv {
	return l.lv
}
