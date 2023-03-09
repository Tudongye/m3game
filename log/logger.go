package log

import (
	"fmt"
	golog "log"
)

func newDefaultLogger() *defaultLogger {
	return &defaultLogger{
		lv: LogLvDebug,
	}
}

type defaultLogger struct {
	lv LogLv
}

func (l *defaultLogger) Output(depth Depth, lv LogLv, plus LogPlus, format string, v ...interface{}) {
	file, funcName, line := depth.String()
	if len(plus) == 0 {
		golog.Printf("[%s] [%s:%d %s] %s\n", lv.String(), file, line, funcName, fmt.Sprintf(format, v...))
	} else {
		golog.Printf("[%s] [%s:%d %s][%s] %s\n", lv.String(), file, line, funcName, plus.String(), fmt.Sprintf(format, v...))
	}
}

func (l *defaultLogger) SetLevel(level LogLv) {
	l.lv = level
}

func (l *defaultLogger) GetLevel() LogLv {
	return l.lv
}
