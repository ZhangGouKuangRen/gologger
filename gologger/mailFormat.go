package gologger

import (
	"fmt"
	"time"
)

type mailFormat struct {
	attributes map[string]string
}

func DefaultMailFormat()*mailFormat  {
	nowTime := time.Now().Format("15:04:05")
	nowDate := time.Now().Format("2006-01-02")
	attributes := make(map[string]string)
	attributes["time"]=nowTime
	attributes["date"]=nowDate
	return &mailFormat{
		attributes: attributes,
	}
}

func (mf *mailFormat)fmtMailMsg(msg *logMsg)string  {
	parsedMsg := `<table border=1 cellspacing=0 width=100% style='text-align:center'>
                    <tr>
                        <th>Date</th>
                        <th>Time</th>
                        <th>File</th>
                        <th>Func</th>
                        <th>Line</th>
                        <th>Level</th>
                        <th>Log</th>
                    </tr>
                    <tr>
                        <td>`+mf.attributes["date"]+`
                        <td>`+mf.attributes["time"]+`
                        <td>`+msg.fileName+`
                        <td>`+msg.funcName+`
                        <td>`+fmt.Sprint(msg.line)+`
                        <td>`+msg.levStr+`
                        <td>`+msg.msg+`
                    </tr>
                </table`
	return parsedMsg
}
