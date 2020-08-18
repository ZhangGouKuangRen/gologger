package gologger

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type gologger struct {
	GlobalLevel       string            `xml:"globalLevel,attr"`
	ConsoleXML        consoleXML        `xml:"console"`
	Wholefiles        wholefiles        `xml:"wholefiles"`
	DailyRollingFiles dailyrollingfiles `xml:"dailyrollingfiles"`
	SizeRollingFiles  sizerollingfiles  `xml:"sizerollingfiles"`
	Smtplog           smtplog           `xml:"smtplog"`
	DatabaseLog       databaselog       `xml:"databaselog"`
	FormatXMLs        formatXMLs        `xml:"formats"`
}
type consoleXML struct {
	Enable    string    `xml:"enable"`
	SelfLevel string    `xml:"selflevel"`
	Format    formatXML `xml:"format"`
}

type formatXMLs struct {
	FormatXML []formatXML `xml:"format"`
}
type formatXML struct {
	Ref    string `xml:"ref,attr"`
	Id     string `xml:"id,attr"`
	Format string `xml:"format,attr"`
}

type wholefiles struct {
	WholeFile []wholefile `xml:"wholefile"`
}
type wholefile struct {
	Path      string    `xml:"path"`
	SelfLevel string    `xml:"selflevel"`
	FormatId  formatXML `xml:"format"`
}

type dailyrollingfiles struct {
	DailyRollingFile []dailyrollingfile `xml:"dailyrollingfile"`
}
type dailyrollingfile struct {
	Path      string    `xml:"path"`
	SelfLevel string    `xml:"selflevel"`
	FormatId  formatXML `xml:"format"`
	Day       string    `xml:"day"`
}

type sizerollingfiles struct {
	SizeRollingFile []sizerollingfile `xml:"sizerollingfile"`
}
type sizerollingfile struct {
	Path      string    `xml:"path"`
	SelfLevel string    `xml:"selflevel"`
	FormatId  formatXML `xml:"format"`
	Size      string    `xml:"size"`
}

type databaselog struct {
	Driver       string `xml:"driver"`
	Username     string `xml:"username"`
	Password     string `xml:"password"`
	Ip           string `xml:"ip"`
	Port         string `xml:"port"`
	Databasename string `xml:"databasename"`
	Tables       tables `xml:"tables"`
}
type tables struct {
	Table []table `xml:"table"`
}
type table struct {
	Name       string     `xml:"name,attr"`
	SelfLevel  string     `xml:"selflevel"`
	Attributes attributes `xml:"attributes"`
}
type attributes struct {
	Attribute []attribute `xml:"attribute"`
}
type attribute struct {
	Name string `xml:"name,attr"`
}

type smtplog struct {
	Host       string     `xml:"host"`
	Port       int64      `xml:"port"`
	Password   string     `xml:"password"`
	From       string     `xml:"from"`
	Subject    string     `xml:"subject"`
	SelfLevel  string     `xml:"selflevel"`
	Recipients recipients `xml:"recipients"`
	Attributes attributes `xml:"attributes"`
}
type recipients struct {
	Recipient []string `xml:"recipient"`
}

func GetLoggerByXML(xmlconfig string) (*logger, error) {
	data, dataErr := ioutil.ReadFile(xmlconfig)
	if dataErr != nil {
		return nil, dataErr
	}
	gologgerXML := gologger{}
	xml.Unmarshal(data, &gologgerXML)

	consoleXML := gologgerXML.ConsoleXML
	wholefiles := gologgerXML.Wholefiles
	dailyrollingfiles := gologgerXML.DailyRollingFiles
	sizerollingfiles := gologgerXML.SizeRollingFiles
	databaselog := gologgerXML.DatabaseLog
	smtplog := gologgerXML.Smtplog

	fmt.Println("consoleXML:",consoleXML)
	fmt.Println("wholefiles:",wholefiles)
	fmt.Println("dailyrollingfiles:", dailyrollingfiles)
	fmt.Println("sizerollingfiles:", sizerollingfiles)
	fmt.Println("databaselog:", databaselog)
	fmt.Println("smtplog:", smtplog)

	formatXMLs := gologgerXML.FormatXMLs
	formatMap := make(map[string]formatXML)
	for _, formatXML := range formatXMLs.FormatXML {
		formatMap[formatXML.Id] = formatXML
	}

	logger := NewLogger(gologgerXML.GlobalLevel)

	consoleSelfLevel, consoleSelfLevelErr := parseLogLevel(consoleXML.SelfLevel)
	if consoleSelfLevelErr != nil {
		return nil, consoleSelfLevelErr
	}
	logger.SetConsoleSelfLogLevel(consoleSelfLevel)
	enableConsole, enconErr := strconv.ParseBool(consoleXML.Enable)
	if enconErr != nil {
		return nil, enconErr
	}
	logger.EnableConsole(enableConsole)

	for _, wholefile := range wholefiles.WholeFile {
		wf := NewWholeFile(wholefile.Path)
		wholefileSelfLevel, wfsErr := parseLogLevel(wholefile.SelfLevel)
		if wfsErr != nil {
			return nil, wfsErr
		}
		wf.SetWholeFileSelfLogLevel(wholefileSelfLevel)
		wf.SetWholeFileFormat(NewFormat(formatMap[wholefile.FormatId.Ref].Format))
		logger.AddWholeFile(wf)
	}

	for _, dailyrollingfile := range dailyrollingfiles.DailyRollingFile {
		drf := NewDailyRollingFile(dailyrollingfile.Path)
		dailyrollingfileSelfLevel, drfLevelErr := parseLogLevel(dailyrollingfile.SelfLevel)
		if drfLevelErr != nil {
			return nil, drfLevelErr
		}
		drf.SetDailyRollingFileSelfLogLevel(dailyrollingfileSelfLevel)
		days, dayErr := strconv.ParseInt(dailyrollingfile.Day, 10, 64)
		if dayErr != nil {
			return nil, dayErr
		}
		drf.SetLogFiletMaxDays(int(days))
		drf.SetDailyRollingFileFormat(NewFormat(formatMap[dailyrollingfile.FormatId.Ref].Format))
		logger.AddDailyRollingFile(drf)
	}

	for _, sizerollingfile := range sizerollingfiles.SizeRollingFile {
		srf := NewSizeRollingFile(sizerollingfile.Path)
		srfLevel, srfLevelErr := parseLogLevel(sizerollingfile.SelfLevel)
		if srfLevelErr != nil {
			return nil, errors.New("sizerollingfile selflevel parse error")
		}
		srf.SetSizeRollingFileSelfLogLevel(srfLevel)
		size, sizeErr := strconv.ParseInt(sizerollingfile.Size, 10, 64)
		if sizeErr != nil {
			return nil, errors.New("sizerollingfile size parse error")
		}
		if size < 0 {
			return nil, errors.New("sizerollingfile size cannot <0")
		}
		srf.SetMaxSize(size * 1024 * 1024)
		srf.SetSizeRollingFileFormat(NewFormat(formatMap[sizerollingfile.FormatId.Ref].Format))
		logger.AddSizeRollingFile(srf)
	}

	port, portErr := strconv.ParseInt(databaselog.Port, 10, 64)
	if portErr != nil {
		return nil, errors.New("xxx/gologger/databaselog/port: cannot parse port to int64")
	}
	ip := net.ParseIP(databaselog.Ip)
	if ip == nil {
		return nil, errors.New("xxx/gologger/databaselog/ip: cannot parse ip")
	}
	dbl := NewDatabaseLog(databaselog.Driver, databaselog.Username, databaselog.Password, ip.String(), port, databaselog.Databasename)
	tables := databaselog.Tables.Table
	for _, table := range tables {
		cols := parseAttributes(table.Attributes.Attribute)
		logTable := newLogTableXML(table.Name, cols)
		logTableSelfLevel, ltslErr := parseLogLevel(table.SelfLevel)
		if ltslErr != nil {
			return nil, errors.New("cannot parse table selflevel")
		}
		logTable.SetTableSelfLevel(logTableSelfLevel)
		dbl.AddTable(logTable)
	}
	logger.AddDatabaseLog(dbl)

	if !verifyEmailFormat(smtplog.From){
		return nil, errors.New("email 'from' format error")
	}
	for _, recpt := range smtplog.Recipients.Recipient {
		if !verifyEmailFormat(recpt){
			return nil, errors.New("email 'recipient' format error")
		}
	}
	sl := NewSmtpLog(smtplog.Host, smtplog.Password, smtplog.From, smtplog.Port, smtplog.Subject)
	smtpLogSelfLevel, smtpLevelErr := parseLogLevel(smtplog.SelfLevel)
	if smtpLevelErr != nil {
		return nil, errors.New("cannot parse smtolog selflevel")
	}
	sl.SetMailSelfLogLevel(smtpLogSelfLevel)
	sl.SetRecipient(smtplog.Recipients.Recipient)
	//logger.AddSmtpLog(sl)

	return logger, nil
}

func parseAttributes(attributes []attribute) []LogAttrubutes {
	cols := []LogAttrubutes{}
	for _, attribute := range attributes {
		switch strings.ToLower(attribute.Name) {
		case "date":
			cols = append(cols, Date)
		case "time":
			cols = append(cols, Time)
		case "file":
			cols = append(cols, File)
		case "func":
			cols = append(cols, Func)
		case "level":
			cols = append(cols, Level)
		case "line":
			cols = append(cols, Line)
		case "msg":
			cols = append(cols, Msg)
		}
	}
	return cols
}

func verifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//func configConsole(consoleXML consoleXML, gologger *gologger)(gologger, error)  {
//	consoleSelfLevel, consoleSelfLevelErr := parseLogLevel(consoleXML.SelfLevel)
//	if consoleSelfLevelErr != nil {
//		return nil, consoleSelfLevelErr
//	}
//	logger.SetConsoleSelfLogLevel(consoleSelfLevel)
//	enableConsole, enconErr := strconv.ParseBool(consoleXML.Enable)
//	if enconErr != nil {
//		return nil, enconErr
//	}
//	logger.EnableConsole(enableConsole)
//	return logger, nil
//}