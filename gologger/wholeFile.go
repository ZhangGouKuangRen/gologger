package gologger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)


type wholeFile struct {
	ofile *os.File
	*format
	hasfmt bool
	fileSelfLev logLevel
	hasFileSelLev bool
	fullpath string
}

func NewWholeFile(fullpath string)*wholeFile {
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
	return &wholeFile{
		ofile: fileptr,
		fullpath: fullpath,
	}
}

func (f *wholeFile) outWholeFile(logMsg *logMsg)  {
	var flag bool
	if f.hasFileSelLev {
		logLev, err := parseLogLevel(logMsg.levStr)
		if err != nil {
			panic("文件"+filepath.Base(f.ofile.Name())+"的私有日志级别不存在")
			return
		}
		if logLev >= f.fileSelfLev {
			flag = true
		}
	} else {
		flag = true
	}
	if flag {
		if f.hasfmt {
			fmt.Fprint(f.ofile, f.fmtMsg(logMsg))
		}else {
			fmt.Fprint(f.ofile, logMsg.msg+"\n")
		}
	}
}

func (f *wholeFile)SetWholeFileFormat(fmt *format)  {
	f.format = fmt
	f.hasfmt = true
}

func (f *wholeFile)SetWholeFileSelfLogLevel(selfLogLev logLevel)  {
	f.fileSelfLev = selfLogLev
	f.hasFileSelLev = true
}