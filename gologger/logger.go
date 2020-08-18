package gologger

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type logLevel uint8

const (
	UNKNOW logLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)
type logger struct {
    level logLevel
    wholeFiles []*wholeFile
    sizeRollingFiles []*sizeRollingFile
    dailyRollingFiles []*dailyRollingFile
    *smtpLog
    *databaseLog
    useDatabaseLog bool
    enconsole bool
    useSmtp bool
    logMsgChan chan *logMsg
    console
    wg *sync.WaitGroup
    mutex *sync.Mutex
	goroutineNum int
}

type logMsg struct {
	msg string
	levStr string
	funcName string
	fileName string
	line int
}

var csl console

//初始化一个logger对象
func NewLogger(levelMsg string)*logger  {
	grteNum := 1
	level, err := parseLogLevel(levelMsg)
	if err != nil {
		panic(err)
	}
	lgr := logger{
		level: level,
		enconsole: true,
		logMsgChan: make(chan *logMsg, 1000),
		//logMsgChan: make(chan *logMsg),
		wg: &sync.WaitGroup{},
		mutex: &sync.Mutex{},
		goroutineNum: grteNum,
	}
	//原本打算开多个goroutine并发写日志，但是考虑到并发写日志时会对共享数据进行修改导致出错，如果加锁解决，效率并不会有所提升
	//所以改用单goroutine写日志(grteNum = 1)
	for i:=0; i<grteNum; i++ {
		go lgr.backgroundWriteLog()
	}
	return &lgr
}

func (lgr *logger)EnableConsole(b bool){
	lgr.enconsole = b
}

func (lgr *logger)backgroundWriteLog()  {
	//lgr.wg.Add(1)
	lgr.mutex.Lock()
	var getMsg bool
	var logMsg *logMsg
	for {
		select {
		case logMsg = <- lgr.logMsgChan:
			getMsg = true
		default:
			time.Sleep(100*time.Millisecond)
			getMsg = false
		}
		if getMsg {
			if lgr.useSmtp {
				go lgr.smtpLog.sendMail(logMsg)
			}
			for _, file := range lgr.wholeFiles {
				file.outWholeFile(logMsg)
			}
			for _, sizeRollingFile := range lgr.sizeRollingFiles {
				sizeRollingFile.outSizeRollingFile(logMsg)
			}
			for _, dailyRollingFile := range lgr.dailyRollingFiles{
				dailyRollingFile.outDailyRollingFile(logMsg)
			}
			if lgr.enconsole {
				lgr.outConsole(logMsg)
			}
			if lgr.useDatabaseLog {
				lgr.insertLog(logMsg)
			}
			lgr.wg.Done()
		}
	}
	lgr.mutex.Unlock()
	//lgr.wg.Done()

}

//格式化记录
func (lgr *logger)log(levStr, logmsg string)  {
	pc, fileName, line, _ := runtime.Caller(3)
	funcName := runtime.FuncForPC(pc).Name()
	logMsg := &logMsg{
		levStr: levStr,
		msg: logmsg,
		fileName: fileName,
		funcName: funcName,
		line: line,
	}
	//放弃使用select，因为lgr.logMsgChan <- logMsg 如果在lgr.wg.Add(1)之前，由于并发原因，在向logMsgChan写入数据后，
	//读出程序可能在wg.add(1)之前立刻读出其中的数据并执行wg.Done();这时可能导致panic: sync: negative WaitGroup counter
	//select {
	//case lgr.logMsgChan <- logMsg:
	//	lgr.wg.Add(1)
	//	fmt.Println("======== +msg ========")
	//default:
	//}
	lgr.wg.Add(1)
	lgr.logMsgChan <- logMsg

}

//生成一条格式化日志
func (lgr *logger)createLog(levStr, msg string, a ...interface{})string  {
	msg = fmt.Sprintf(msg, a...)
	return msg
}

func (lgr *logger)AddWholeFile(fil *wholeFile)  {
	lgr.wholeFiles = append(lgr.wholeFiles, fil)
}
func (lgr *logger)AddSizeRollingFile(rf *sizeRollingFile)  {
	lgr.sizeRollingFiles = append(lgr.sizeRollingFiles, rf)
}

func (lgr *logger)AddDailyRollingFile(rf *dailyRollingFile)  {
	lgr.dailyRollingFiles = append(lgr.dailyRollingFiles, rf)
}

func (lgr *logger)AddSmtpLog(smtplog *smtpLog)  {
	lgr.useSmtp = true
	lgr.smtpLog = smtplog
}

func (lgr *logger)AddDatabaseLog(databaseLog *databaseLog)  {
	lgr.useDatabaseLog = true
	lgr.databaseLog = databaseLog
}

func (lgr *logger)Flush(){
	lgr.wg.Wait()
	for _, wholeFile := range lgr.wholeFiles {
		wholeFile.ofile.Close()
	}
	for _, sizeRollingFile := range lgr.sizeRollingFiles {
		sizeRollingFile.ofile.Close()
	}
	for _, dailyRollingFile := range lgr.dailyRollingFiles {
		dailyRollingFile.ofile.Close()
	}
	if lgr.databaseLog != nil {
		lgr.databaseLog.db.Close()
	}
}
