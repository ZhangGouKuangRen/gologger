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
func (lgr *logger) enableLog(level logLevel) bool {
	if lgr.level <= level {
		return true
	} else {
		return false
	}
}

func (lgr *logger) parseMsgf(levStr, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	lgr.log(levStr, msg)
}
func (lgr *logger) parseMsg(levStr string, a ...interface{}) {
	msg := fmt.Sprint(a...)
	lgr.log(levStr, msg)
}

//Debug级别日志
func (lgr *logger) Debugf(format string, a ...interface{}) {
	if lgr.enableLog(DEBUG) {
		lgr.parseMsgf("Debug", format, a...)
	}
}

func (lgr *logger) Tracef(format string, a ...interface{}) {
	if lgr.enableLog(TRACE) {
		lgr.parseMsgf("Trace", format, a...)
	}
}

func (lgr *logger) Infof(format string, a ...interface{}) {
	if lgr.enableLog(INFO) {
		lgr.parseMsgf("Info", format, a...)
	}
}

func (lgr *logger) Warningf(format string, a ...interface{}) {
	if lgr.enableLog(WARNING) {
		lgr.parseMsgf("Warning", format, a...)
	}
}

func (lgr *logger) Errorf(format string, a ...interface{}) {
	if lgr.enableLog(ERROR) {
		lgr.parseMsgf("Error", format, a...)
	}
}

func (lgr *logger) Fatalf(format string, a ...interface{}) {
	if lgr.enableLog(FATAL) {
		lgr.parseMsgf("Fatal", format, a...)
	}
}

func (lgr *logger) Debug(a ...interface{}) {
	if lgr.enableLog(DEBUG) {
		lgr.parseMsg("Debug", a...)
	}
}

func (lgr *logger) Trace(a ...interface{}) {
	if lgr.enableLog(TRACE) {
		lgr.parseMsg("Trace", a...)
	}
}

func (lgr *logger) Info(a ...interface{}) {
	if lgr.enableLog(INFO) {
		lgr.parseMsg("Info", a...)
	}
}

func (lgr *logger) Warning(a ...interface{}) {
	if lgr.enableLog(WARNING) {
		lgr.parseMsg("Warning", a...)
	}
}

func (lgr *logger) Error(a ...interface{}) {
	if lgr.enableLog(ERROR) {
		lgr.parseMsg("Error", a...)
	}
}

func (lgr *logger) Fatal(a ...interface{}) {
	if lgr.enableLog(FATAL) {
		lgr.parseMsg("Fatal", a...)
	}
}
