package log

import (
	"fmt"
	goruntime "runtime"
	"strings"
)

var (
	_logger        Logger
	_defaultlogger = newDefaultLogger()
)

const (
	LogLvDebug LogLv = iota + 1
	LogLvInfo
	LogLvWarn
	LogLvError
	LogLvFatal
)

const (
	LogFlagDebug = "DEBUG"
	LogFlagInfo  = "INFO "
	LogFlagWarn  = "WARN "
	LogFlagError = "ERROR"
	LogFlagFatal = "FATAL"
)

const (
	_outputdepth Depth = 2
)

type LogLv int

func (l LogLv) String() string {
	switch l {
	case LogLvDebug:
		return LogFlagDebug
	case LogLvInfo:
		return LogFlagInfo
	case LogLvWarn:
		return LogFlagWarn
	case LogLvError:
		return LogFlagError
	case LogLvFatal:
		return LogFlagFatal
	default:
		return "unknow"
	}
}

func ConvertLogLv(s string) LogLv {
	switch s {
	case LogFlagDebug:
		return LogLvDebug
	case LogFlagInfo:
		return LogLvInfo
	case LogFlagWarn:
		return LogLvWarn
	case LogFlagError:
		return LogLvError
	case LogFlagFatal:
		return LogLvFatal
	default:
		panic(fmt.Sprintf("Unknow LogLv %s", s))
	}
}

func Get() Logger {
	if _logger == nil {
		return _defaultlogger
	}
	return _logger
}

func Set(logger Logger) {
	if _logger != nil {
		panic("")
	}
	_logger = logger
}

type LogPlus map[string]string

func (l *LogPlus) String() string {
	var s []string
	for k, v := range *l {
		s = append(s, fmt.Sprintf("%s:%s", k, v))
	}
	ps := strings.Join(s, "")
	return ps
}

type Depth int

func (d *Depth) String() (string, string, int) {
	const callOffset = 1
	pc, file, line, ok := goruntime.Caller(int(*d) + callOffset)
	if !ok {
		return "", "", 0
	}
	funcName := goruntime.FuncForPC(pc).Name()
	idx := strings.LastIndexByte(funcName, '.')
	if idx != -1 {
		funcName = funcName[idx+1:]
	}

	idx = strings.LastIndexByte(file, '/')
	if idx == -1 {
		return file, funcName, line
	}
	idx = strings.LastIndexByte(file[:idx], '/')
	if idx == -1 {
		return file, funcName, line
	}
	return file[idx+1:], funcName, line
}

type Logger interface {
	Output(depth Depth, lv LogLv, plus LogPlus, format string, v ...interface{})
	SetLevel(level LogLv)
	GetLevel() LogLv
}

func Info(format string, v ...interface{}) {
	if Get().GetLevel() > LogLvInfo {
		return
	}
	Get().Output(_outputdepth, LogLvInfo, nil, format, v...)
}

func Debug(format string, v ...interface{}) {
	if Get().GetLevel() > LogLvDebug {
		return
	}
	Get().Output(_outputdepth, LogLvDebug, nil, format, v...)
}

func Warn(format string, v ...interface{}) {
	if Get().GetLevel() > LogLvWarn {
		return
	}
	Get().Output(_outputdepth, LogLvWarn, nil, format, v...)
}

func Error(format string, v ...interface{}) {
	if Get().GetLevel() > LogLvError {
		return
	}
	Get().Output(_outputdepth, LogLvError, nil, format, v...)
}

func Fatal(format string, v ...interface{}) {
	if Get().GetLevel() > LogLvFatal {
		return
	}
	Get().Output(_outputdepth, LogLvFatal, nil, format, v...)
}

func InfoP(plus LogPlus, format string, v ...interface{}) {
	if Get().GetLevel() > LogLvInfo {
		return
	}
	Get().Output(_outputdepth, LogLvInfo, plus, format, v...)
}

func DebugP(plus LogPlus, format string, v ...interface{}) {
	if Get().GetLevel() > LogLvDebug {
		return
	}
	Get().Output(_outputdepth, LogLvDebug, plus, format, v...)
}

func WarnP(plus LogPlus, format string, v ...interface{}) {
	if Get().GetLevel() > LogLvWarn {
		return
	}
	Get().Output(_outputdepth, LogLvWarn, plus, format, v...)
}

func ErrorP(plus LogPlus, format string, v ...interface{}) {
	if Get().GetLevel() > LogLvError {
		return
	}
	Get().Output(_outputdepth, LogLvError, plus, format, v...)
}

func FatalP(plus LogPlus, format string, v ...interface{}) {
	if Get().GetLevel() > LogLvFatal {
		return
	}
	Get().Output(_outputdepth, LogLvFatal, plus, format, v...)
}
