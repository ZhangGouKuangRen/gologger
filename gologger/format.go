package gologger

import (
	"fmt"
	"strconv"
	"time"
)

type format struct {
	parsedFormat string
	originalFormat string
	isDeaultFormat bool
}

var flatWords = map[string]string{
    "Level": "Level",
    "Date" : "",
}

func DefaultFormat()*format {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	ffmt := fmt.Sprintf("[%s]", nowTime)+"[%s][%s：%s：%d]"+"%s\n"
	return &format{
		parsedFormat: ffmt,
		isDeaultFormat: true,
	}
}

func NewFormat(originalFormat string)*format {
	format := &format{
		originalFormat: originalFormat,
		isDeaultFormat: false,
	}
	return format
}

func (frmt *format)fmtMsg(logMsg *logMsg)string  {
    var log string	
 	if frmt.isDeaultFormat {
		log = fmt.Sprintf(frmt.parsedFormat, logMsg.levStr, logMsg.fileName, logMsg.funcName, logMsg.line, logMsg.msg)
	} else {
		log = frmt.parseFormat(logMsg)
	}
	return log
}

func (frmt *format)parseFormat(logMsg *logMsg)string {
	indexMap := make(map[int]string)
	keywordsValue := []string{}
	for i:=0; i<len(frmt.originalFormat); i++{
		keyword := []rune{}
		c := frmt.originalFormat[i]
		if c == '%' {
			i++
			temp := i
			for ; i<len(frmt.originalFormat); i++{
				if string(frmt.originalFormat[i]) == "%"{
					i--
					break
				}
				keyword = append(keyword, rune(frmt.originalFormat[i]))
				result := getParsedValue(string(keyword), logMsg)
				if result != "" {
					indexMap[temp-1]=strconv.FormatInt(int64(i+1), 10)
					keywordsValue = append(keywordsValue, result)
					break
				}
			}
		}
	}
	parsedFormat := ""
	valueIndex := 0
	for indexStart:=0; indexStart<len(frmt.originalFormat); indexStart++{
		indexEnd := indexMap[indexStart]
		if indexEnd != "" {
			end, _ :=strconv.ParseInt(indexEnd, 10, 64)
			value := keywordsValue[valueIndex]
			valueIndex++
			parsedFormat = parsedFormat+value
			indexStart=int(end-1)
		} else {
			parsedFormat+=string(frmt.originalFormat[indexStart])
		}
	}

	return parsedFormat+"\n"
}

func getDate()string  {
	dateStr := time.Now().Format("2006-01-02")
	return dateStr
}

func getTime()string  {
	timeStr := time.Now().Format("15:04:05")
	return timeStr
}

func getParsedValue(keyword string, logMsg *logMsg)string  {
	switch keyword {
	case "Date":
		return getDate()
	case "Time":
		return getTime()
	case "Level":
		return logMsg.levStr
	case "File":
		return logMsg.fileName
	case "Func":
		return logMsg.funcName
	case "Line":
		return strconv.FormatInt(int64(logMsg.line), 10)
	case "Msg":
		return logMsg.msg
	default:
		return ""
	}
}