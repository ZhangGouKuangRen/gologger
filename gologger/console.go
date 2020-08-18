package gologger

import (
	"fmt"
	"os"
	"strings"
)

type console struct {
	*format
	hasfmt bool
	conSelfLev logLevel
	hasConSelLev bool
}

func (csl *console)outConsole(logMsg *logMsg)  {
	var flag bool
	if csl.hasConSelLev {
		logLev, err := parseLogLevel(logMsg.levStr)
		if err != nil {
			panic("控制台私有日志级别不存在")
			return
		}
		if logLev >= csl.conSelfLev {
			flag = true
		}
	} else {
		flag = true
	}
	if flag {
		if csl.hasfmt {
			fmt.Fprint(os.Stdout, csl.fmtMsg(logMsg))
		}else {
			fmt.Fprintln(os.Stdout, "["+strings.ToUpper(logMsg.levStr)+"] "+logMsg.msg)
		}
	}
}

func (csl *console)SetConsoleFormat(fmt *format)  {
	csl.format = fmt
	csl.hasfmt = true
}

func (csl *console)SetConsoleSelfLogLevel(selfLogLev logLevel)  {
	csl.conSelfLev = selfLogLev
	csl.hasConSelLev = true
}