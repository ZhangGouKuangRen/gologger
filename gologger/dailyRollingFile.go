package gologger

import (
	"errors"
	"fmt"
	"gologger/util"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type dailyRollingFile struct {
	currentDate string
	ofile *os.File
	*format
	hasfmt bool
	fileSelfLev logLevel
	hasFileSelLev bool
	fullpath string
	maxDays int
	//这里使用队列
	logFileNames *util.Queue
}

func NewDailyRollingFile(fullpath string)*dailyRollingFile  {
	var fileptr *os.File
	queue := util.NewQueue()

	currentDate := time.Now().Format("20060102")
	f, err := os.OpenFile(fullpath, os.O_RDWR | os.O_CREATE | os.O_APPEND, os.ModeAppend|os.ModePerm)
	index := strings.LastIndex(fullpath, string("/"))
	path := string([]rune(fullpath)[:index])
	if err != nil {
		if os.IsNotExist(err) {
			mkerr := os.MkdirAll(path, os.ModeAppend)
			if mkerr != nil {
				panic("创建文件目录失败")
			}
			fc, ferr := os.Create(fullpath)
			if ferr != nil {
				panic("创建文件失败")
			}
			fileptr = fc
		} else {
			panic(err)
		}
	} else {
		fileptr = f
	}

	return &dailyRollingFile{
		ofile: fileptr,
		currentDate: currentDate,
		fullpath: fullpath,
		maxDays: 7,
		logFileNames: queue,

	}
}


func (drf *dailyRollingFile) outDailyRollingFile(logMsg *logMsg)  {
	var flag bool
	if drf.hasFileSelLev {
		logLev, err := parseLogLevel(logMsg.levStr)
		if err != nil {
			panic("文件"+filepath.Base(drf.ofile.Name())+"的私有日志级别不存在")
			return
		}
		if logLev >= drf.fileSelfLev {
			flag = true
		}
	} else {
		flag = true
	}
	if flag {
		if drf.hasfmt {
			drf.cutDailyRollingFile(logMsg.levStr, drf.fmtMsg(logMsg))
			fmt.Fprint(drf.ofile, drf.fmtMsg(logMsg))
		}else {
			drf.cutDailyRollingFile(logMsg.levStr, logMsg.msg)
			fmt.Fprint(drf.ofile, logMsg.msg+"\n")
		}
	}
}

func (drf *dailyRollingFile)SetDailyRollingFileFormat(fmt *format)  {
	drf.format = fmt
	drf.hasfmt = true
}

func (drf *dailyRollingFile)SetDailyRollingFileSelfLogLevel(selfLogLev logLevel)  {
	drf.fileSelfLev = selfLogLev
	drf.hasFileSelLev = true
}

func (drf *dailyRollingFile)SetLogFiletMaxDays(newMaxDays int)  {
	if newMaxDays < 0 {
		panic(errors.New("newMaxDays不能为负数"))
	}
	drf.maxDays = newMaxDays
}

//如果当前日期与currentDate不相同，则保存文件，更新currentDate为当前日期，并创建一个新的文件；当文件数大于maxDays时，删除最旧的文件然后创建新文件；如果maxDays为0，则不删除旧文件
func (drf *dailyRollingFile) cutDailyRollingFile(levStr, msg string){
	if drf.maxDays < drf.logFileNames.Size() && drf.maxDays != 0{
		//logFileNames出队
		logFileName, pollLogFileErr := drf.logFileNames.Poll()
		if pollLogFileErr != nil {
			panic(pollLogFileErr)
		}
		//删除一个最旧文件logFileName
		err := os.Remove(fmt.Sprint(logFileName))
		if err != nil {
			panic(err)
		}
	}
	fileInfo, fierr := drf.ofile.Stat()
	if fierr != nil {
		panic(fierr)
	}
    cuntDate := time.Now().Format("20060102")
	if cuntDate != drf.currentDate {
		fname := filepath.Base(fileInfo.Name())
		newName := getDailyRollingNewName(drf.fullpath, drf.currentDate, fname)
		drf.ofile.Close()
		rerr := os.Rename(drf.fullpath,  newName)
		//将生成的新log文件名保存到logFileNames队列
		drf.logFileNames.Push(newName)
		if rerr != nil{
			panic(rerr)
		}
		nfile, ferr :=os.OpenFile(drf.fullpath, os.O_CREATE|os.O_APPEND, 0600)
		if ferr != nil {
			panic(ferr)
		}
		drf.ofile = nfile
		drf.currentDate = cuntDate
	}
}

func getDailyRollingNewName(fpath,currentDate, oldname string)string{
	on := string(fpath[:len(fpath)-len(oldname)])
	var prefix, suffix string
	index := strings.LastIndex(oldname, ".")
	if index == -1 {
		prefix = oldname
		suffix = ""
	} else {
		prefix = string([]rune(oldname)[:index])
		suffix = string([]rune(oldname)[index:])
	}
	mid := currentDate
	newname := prefix+mid+suffix
	return on+newname
}