package log

import (
	goruntime "runtime"
	"strings"
)

var (
	_logger Logger
)

const (
	LogLvNull LogLv = iota
	LogLvTrace
	LogLvDebug
	LogLvInfo
	LogLvWarn
	LogLvError
	LogLvFatal
)

const (
	LogFlagNull  = "NULL "
	LogFlagTrace = "TRACE"
	LogFlagDebug = "DEBUG"
	LogFlagInfo  = "INFO "
	LogFlagWarn  = "WARN "
	LogFlagError = "ERROR"
	LogFlagFatal = "FATAL"
)

const (
	_outputdepth = 2
)

func init() {
	_logger = newDefaultLogger()
}

type LogLv int

func (l LogLv) String() string {
	switch l {
	case LogLvTrace:
		return LogFlagTrace
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
		return LogFlagNull
	}
}

func GetLogger() Logger {
	return _logger
}

func SerLogger(logger Logger) {
	_logger = logger
}

type Logger interface {
	Output(depth int, logflag string, plus string, format string, v ...interface{})
	SetLevel(level LogLv)
	GetLevel() LogLv
}

func Trace(format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvTrace {
		return
	}
	_logger.Output(_outputdepth, LogFlagTrace, "", format, v...)
}

func Info(format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvInfo {
		return
	}
	_logger.Output(_outputdepth, LogFlagInfo, "", format, v...)
}

func Debug(format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvDebug {
		return
	}
	_logger.Output(_outputdepth, LogFlagDebug, "", format, v...)
}

func Warn(format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvWarn {
		return
	}
	_logger.Output(_outputdepth, LogFlagWarn, "", format, v...)
}

func Error(format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvError {
		return
	}
	_logger.Output(_outputdepth, LogFlagError, "", format, v...)
}

func Fatal(format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvFatal {
		return
	}
	_logger.Output(_outputdepth, LogFlagFatal, "", format, v...)
}

func TraceP(plus string, format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvTrace {
		return
	}
	_logger.Output(_outputdepth, LogFlagTrace, plus, format, v...)
}

func InfoP(plus string, format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvInfo {
		return
	}
	_logger.Output(_outputdepth, LogFlagInfo, plus, format, v...)
}

func DebugP(plus string, format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvDebug {
		return
	}
	_logger.Output(_outputdepth, LogFlagDebug, plus, format, v...)
}

func WarnP(plus string, format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvWarn {
		return
	}
	_logger.Output(_outputdepth, LogFlagWarn, plus, format, v...)
}

func ErrorP(plus string, format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvError {
		return
	}
	_logger.Output(_outputdepth, LogFlagError, plus, format, v...)
}

func FatalP(plus string, format string, v ...interface{}) {
	if _logger == nil {
		return
	}
	if _logger.GetLevel() > LogLvFatal {
		return
	}
	_logger.Output(_outputdepth, LogFlagFatal, plus, format, v...)
}

func GetCallPath(depth int) (string, string, int) {
	const callOffset = 1
	pc, file, line, ok := goruntime.Caller(depth + callOffset)
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
