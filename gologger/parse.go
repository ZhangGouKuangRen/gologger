package gologger

import (
	"fmt"
	"strings"
)

//解析日志记录级别
func parseLogLevel(levelMsg string) (logLevel, error) {
	var level logLevel
	switch strings.ToLower(levelMsg) {
	case "debug":
		level = DEBUG
		return level, nil
	case "trace":
		level = TRACE
		return level, nil
	case "info":
		level = INFO
		return level, nil
	case "warning":
		level = WARNING
		return level, nil
	case "error":
		level = ERROR
		return level, nil
	case "fatal":
		level = FATAL
		return level, nil
	default:
		level = UNKNOW
		return UNKNOW, fmt.Errorf("%s", "错误的日志级别")
	}
}

//日志级别比较
func (lgr *Logger) enableLog(level logLevel) bool {
	if lgr.level <= level {
		return true
	} else {
		return false
	}
}

func (lgr *Logger) parseMsgf(levStr, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	lgr.log(levStr, msg)
}
func (lgr *Logger) parseMsg(levStr string, a ...interface{}) {
	msg := fmt.Sprint(a...)
	lgr.log(levStr, msg)
}

//Debug级别日志
func (lgr *Logger) Debugf(format string, a ...interface{}) {
	if lgr.enableLog(DEBUG) {
		lgr.parseMsgf("Debug", format, a...)
	}
}

func (lgr *Logger) Tracef(format string, a ...interface{}) {
	if lgr.enableLog(TRACE) {
		lgr.parseMsgf("Trace", format, a...)
	}
}

func (lgr *Logger) Infof(format string, a ...interface{}) {
	if lgr.enableLog(INFO) {
		lgr.parseMsgf("Info", format, a...)
	}
}

func (lgr *Logger) Warningf(format string, a ...interface{}) {
	if lgr.enableLog(WARNING) {
		lgr.parseMsgf("Warning", format, a...)
	}
}

func (lgr *Logger) Errorf(format string, a ...interface{}) {
	if lgr.enableLog(ERROR) {
		lgr.parseMsgf("Error", format, a...)
	}
}

func (lgr *Logger) Fatalf(format string, a ...interface{}) {
	if lgr.enableLog(FATAL) {
		lgr.parseMsgf("Fatal", format, a...)
	}
}

func (lgr *Logger) Debug(a ...interface{}) {
	if lgr.enableLog(DEBUG) {
		lgr.parseMsg("Debug", a...)
	}
}

func (lgr *Logger) Trace(a ...interface{}) {
	if lgr.enableLog(TRACE) {
		lgr.parseMsg("Trace", a...)
	}
}

func (lgr *Logger) Info(a ...interface{}) {
	if lgr.enableLog(INFO) {
		lgr.parseMsg("Info", a...)
	}
}

func (lgr *Logger) Warning(a ...interface{}) {
	if lgr.enableLog(WARNING) {
		lgr.parseMsg("Warning", a...)
	}
}

func (lgr *Logger) Error(a ...interface{}) {
	if lgr.enableLog(ERROR) {
		lgr.parseMsg("Error", a...)
	}
}

func (lgr *Logger) Fatal(a ...interface{}) {
	if lgr.enableLog(FATAL) {
		lgr.parseMsg("Fatal", a...)
	}
}
