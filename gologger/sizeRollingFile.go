package gologger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type sizeRollingFile struct {
	maxsize int64
	ofile *os.File
	*format
	hasfmt bool
	fileSelfLev logLevel
	hasFileSelLev bool
	fullpath string
}

func NewSizeRollingFile(fullpath string, size ...int64)*sizeRollingFile {
	var fileptr *os.File
	//f, err := os.Create(fullpath)
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
	var s int64
	s = 1024*1024*10
	if len(size) == 1 {
		s = size[0]
	}
	return &sizeRollingFile{
		ofile: fileptr,
		maxsize: s,
		fullpath: fullpath,
	}
}

func (srf *sizeRollingFile) outSizeRollingFile(logMsg *logMsg)  {
	var flag bool
	if srf.hasFileSelLev {
		logLev, err := parseLogLevel(logMsg.levStr)
		if err != nil {
			panic("文件"+filepath.Base(srf.ofile.Name())+"的私有日志级别不存在")
			return
		}
		if logLev >= srf.fileSelfLev {
			flag = true
		}
	} else {
		flag = true
	}
	if flag {
		if srf.hasfmt {
			srf.cutSizeRollingFile(logMsg.levStr, srf.fmtMsg(logMsg))
			fmt.Fprint(srf.ofile, srf.fmtMsg(logMsg))
		}else {
			srf.cutSizeRollingFile(logMsg.levStr, logMsg.msg)
			fmt.Fprint(srf.ofile, logMsg.msg+"\n")
		}
	}
}

func (rf *sizeRollingFile)SetSizeRollingFileFormat(fmt *format)  {
	rf.format = fmt
	rf.hasfmt = true
}

func (rf *sizeRollingFile)SetSizeRollingFileSelfLogLevel(selfLogLev logLevel)  {
	rf.fileSelfLev = selfLogLev
	rf.hasFileSelLev = true
}

func (rf *sizeRollingFile)SetMaxSize(size int64)  {
	rf.maxsize = size
}

//保存达到maxsize的日志，创建一个新的文件
func (rf *sizeRollingFile) cutSizeRollingFile(levStr, msg string){
	fileInfo, fierr :=rf.ofile.Stat()
	if fierr != nil {
		panic(fierr)
	}
	fsize := fileInfo.Size()
	lsize := int64(len([]rune(msg)))
	if fsize+lsize > rf.maxsize {
		fname := filepath.Base(fileInfo.Name())
		newName := getSizeRollingNewName(rf.fullpath, fname)
		rf.ofile.Close()
		rerr := os.Rename(rf.fullpath,  newName)
		if rerr != nil{
			panic(rerr)
		}
		nfile, ferr :=os.OpenFile(rf.fullpath, os.O_CREATE|os.O_APPEND, 0600)
		if ferr != nil {
			panic(ferr)
		}
		rf.ofile = nfile
		}
}

func getSizeRollingNewName(fpath, oldname string)string{
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
	mid := time.Now().Format("20060102_150405000")
	newname := prefix+mid+suffix
	return on+newname
}