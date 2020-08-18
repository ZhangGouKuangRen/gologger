package mailutil

import (
	"fmt"
	"net/smtp"
)
type MailUtil struct {
	*mailConfig
	*messageConfig
}

type mailConfig struct {
	host string
	port int64
	passwd string
	from string
	nickName string
}

type messageConfig struct {
	contentType string
	charset string
}

func NewMailConfig(host, passwd, from, nickName string, port int64)*mailConfig  {
	return &mailConfig{
		host: host,
		port: port,
		passwd: passwd,
		from: from,
		nickName: nickName,
	}
}

func NewMessageConfig(contentType, charset string)*messageConfig  {
	return &messageConfig{
		contentType: contentType,
		charset: charset,
	}
}

func NewMailUtil(mailConfig *mailConfig, messageConfig *messageConfig)*MailUtil  {
	mailUtil := &MailUtil{
		mailConfig,
		messageConfig,
	}
	return mailUtil
}

func (mailUtil *MailUtil)SendMail(subject, body string, to []string)error  {
	if len(to)<=0 {
		return nil
	}
	auth := smtp.PlainAuth("", mailUtil.from, mailUtil.passwd, mailUtil.host)
	msg := []byte("To: " + fmt.Sprint(to) + "\r\nFrom: "+ mailUtil.nickName + "<" + mailUtil.from + ">\r\nSubject: " + subject + "\r\n" + "Content-Type: "+mailUtil.contentType+"; charset="+mailUtil.charset + "\r\n\r\n" + body)
	addr := mailUtil.host+":"+fmt.Sprint(mailUtil.port)
	err := smtp.SendMail(addr, auth, mailUtil.from, to, []byte(msg))
	return err
}

func (msgcfg *messageConfig)SetCharset(charset string)  {

}
func (msgcfg *messageConfig)SetContentType(contentType string)  {

}