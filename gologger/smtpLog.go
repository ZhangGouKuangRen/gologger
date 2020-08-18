package gologger

import (
	"gologger/util/mailutil"
)

type smtpLog struct {
	mailutil *mailutil.MailUtil
	subject string
	to []string
	hasMailSelfLevel bool
	mailSelfLevel logLevel
	mailFormat *mailFormat
}

func NewSmtpLog(host, passwd, from string, port int64, subject string)*smtpLog  {
	mailConfig := mailutil.NewMailConfig(host, passwd, from, "", port)
	messageConfig := mailutil.NewMessageConfig("text/html", "UTF-8")
	mailutil := mailutil.NewMailUtil(mailConfig, messageConfig)
	format := DefaultMailFormat()
    return &smtpLog{
		mailutil:mailutil,
    	subject:subject,
    	to:[]string{},
    	hasMailSelfLevel: true,
    	mailSelfLevel: ERROR,
    	mailFormat: format,
	}
}

func (sl *smtpLog)SetRecipient(recipient []string)  {
	sl.to = recipient
}

func (sl *smtpLog) SetMailSelfLogLevel(level logLevel)  {
	sl.hasMailSelfLevel=true
	sl.mailSelfLevel=level
}

func (sl *smtpLog) sendMail(msg *logMsg) {
	if len(sl.to)<=0 {
		return
	}
	var flag bool
	if sl.hasMailSelfLevel {
		logLev, err := parseLogLevel(msg.levStr)
		if err != nil {
			panic(err)
			return
		}
		if logLev >= sl.mailSelfLevel {
			flag = true
		}
	}else {
		flag = true
	}
	if flag {
		mailErr := sl.mailutil.SendMail(sl.subject, sl.mailFormat.fmtMailMsg(msg), sl.to)
		if mailErr != nil {
			panic(mailErr)
		}
	}

}
