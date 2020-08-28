package gologger

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var formatMap = make(map[string]formatXML)

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

func GetLoggerByXML(xmlconfig string) (*Logger, error) {
	data, dataErr := ioutil.ReadFile(xmlconfig)
	if dataErr != nil {
		return nil, dataErr
	}
	gologgerXML := gologger{}
	xml.Unmarshal(data, &gologgerXML)

	console := gologgerXML.ConsoleXML
	wholefile := gologgerXML.Wholefiles
	dailyrollingfile := gologgerXML.DailyRollingFiles
	sizerollingfile := gologgerXML.SizeRollingFiles
	databaselg := gologgerXML.DatabaseLog
	smtplg := gologgerXML.Smtplog

	formatXMLs := gologgerXML.FormatXMLs
	for _, formatXML := range formatXMLs.FormatXML {
		formatMap[formatXML.Id] = formatXML
	}
	logger := NewLogger(gologgerXML.GlobalLevel)

	if !reflect.DeepEqual(console, consoleXML{}){
		checkErr(configConsole(console, logger))
	}
	if !reflect.DeepEqual(wholefile, wholefiles{}){
		checkErr(configWholeFiles(wholefile, logger))
	}
	if !reflect.DeepEqual(dailyrollingfile, dailyrollingfiles{}){
		checkErr(configDailyRollingFiles(dailyrollingfile, logger))
	}
	if !reflect.DeepEqual(sizerollingfile, sizerollingfiles{}){
		checkErr(configSizeRollingFiles(sizerollingfile, logger))
	}
	if !reflect.DeepEqual(databaselg, databaselog{}){
		checkErr(configDatabaseLog(databaselg, logger))
	}
	if !reflect.DeepEqual(smtplg, smtplog{}){
		checkErr(configSmtpLog(smtplg, logger))
	}
	return logger, nil
}

func parseAttributes(attributes []attribute) ([]LogAttrubutes, error) {
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
		default:
			return nil, errors.New("databaselog table attribute "+attribute.Name+" cannot be parsed to gologger keyword")
		}
	}
	return cols, nil
}

func verifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func configConsole(consoleXML consoleXML, logger *Logger)(*Logger, error)  {
	selfLevel := consoleXML.SelfLevel
	if selfLevel != "" {
		consoleSelfLevel, consoleSelfLevelErr := parseLogLevel(selfLevel)
		if consoleSelfLevelErr != nil {
			return nil, errors.New("console self level cannot be parse")
		}
		logger.SetConsoleSelfLogLevel(consoleSelfLevel)
	}
	enable := consoleXML.Enable
	if enable != "" {
		enableConsole, enconErr := strconv.ParseBool(enable)
		if enconErr != nil {
			return nil, errors.New("console 'enable' value cannot be parsed to true or false")
		}
		logger.EnableConsole(enableConsole)
	}
	format, hasFormat := parseFormatRef(consoleXML.Format.Ref)
	if hasFormat {
		logger.SetConsoleFormat(format)
	}
	return logger, nil
}

func configWholeFiles(wholefiles wholefiles, logger *Logger)(*Logger, error)  {
	for _, wholefile := range wholefiles.WholeFile {
		path := wholefile.Path
		if path != "" {
			wf := NewWholeFile(wholefile.Path)
			selfLevel := wholefile.SelfLevel
			if selfLevel != "" {
				wholefileSelfLevel, wfsErr := parseLogLevel(selfLevel)
				if wfsErr != nil {
					return nil, errors.New("wholefile selflevel cannot be parsed")
				}
				wf.SetWholeFileSelfLogLevel(wholefileSelfLevel)
				format, hasFormat := parseFormatRef(wholefile.FormatId.Ref)
				if hasFormat {
					wf.SetWholeFileFormat(format)
				}
			}
			logger.AddWholeFile(wf)
		}
	}
	return logger, nil
}

func configDailyRollingFiles(dailyrollingfiles dailyrollingfiles, logger *Logger)(*Logger, error)  {
	for _, dailyrollingfile := range dailyrollingfiles.DailyRollingFile {
		path := dailyrollingfile.Path
		if path != "" {
			drf := NewDailyRollingFile(path)
			selfLevelStr := dailyrollingfile.SelfLevel
			if selfLevelStr != "" {
				dailyrollingfileSelfLevel, drfLevelErr := parseLogLevel(selfLevelStr)
				if drfLevelErr != nil {
					return nil, errors.New("wholefile selflevel cannot be parsed")
				}
				drf.SetDailyRollingFileSelfLogLevel(dailyrollingfileSelfLevel)
			}
			daystr := dailyrollingfile.Day
			if daystr != "" {
				days, dayErr := strconv.ParseInt(daystr, 10, 64)
				if dayErr != nil {
					return nil, errors.New("dailyrollingfile day cannot be parsed to int")
				}
				drf.SetLogFiletMaxDays(int(days))
				format, hasFormat := parseFormatRef(dailyrollingfile.FormatId.Ref)
				if hasFormat {
					drf.SetDailyRollingFileFormat(format)
				}
			}
			logger.AddDailyRollingFile(drf)
		}
	}
	return logger, nil
}

func configSizeRollingFiles(sizerollingfiles sizerollingfiles, logger *Logger)(*Logger, error)  {
	for _, sizerollingfile := range sizerollingfiles.SizeRollingFile {
		path := sizerollingfile.Path
		if path != "" {
			srf := NewSizeRollingFile(path)
			selfLevelStr := sizerollingfile.SelfLevel
			if selfLevelStr != "" {
				srfLevel, srfLevelErr := parseLogLevel(selfLevelStr)
				if srfLevelErr != nil {
					return nil, errors.New("sizerollingfile selflevel cannot be parsed")
				}
				srf.SetSizeRollingFileSelfLogLevel(srfLevel)
			}
			sizeStr := sizerollingfile.Size
			if sizeStr != "" {
				size, sizeErr := strconv.ParseInt(sizeStr, 10, 64)
				if sizeErr != nil {
					return nil, errors.New("sizerollingfile size cannot be parsed to int")
				}
				if size < 0 {
					return nil, errors.New("sizerollingfile size cannot <0")
				}
				srf.SetMaxSize(size * 1024 * 1024)
			}
			format, hasFormat := parseFormatRef(sizerollingfile.FormatId.Ref)
			if hasFormat {
				srf.SetSizeRollingFileFormat(format)
			}
			logger.AddSizeRollingFile(srf)
		}
	}
	return logger, nil
}

func configDatabaseLog(databaselog databaselog, logger *Logger)(*Logger, error)  {
	port, portErr := strconv.ParseInt(databaselog.Port, 10, 64)
	if portErr != nil {
		return nil, errors.New("databaselog port cannot be parsed to int64")
	}
	ip := net.ParseIP(databaselog.Ip)
	if ip == nil {
		return nil, errors.New("databaselog ip cannot be parsed to IP")
	}
	driver := databaselog.Driver
	username := databaselog.Username
	password := databaselog.Password
	databaseName := databaselog.Databasename
	switch "" {
	case driver:
		return nil, errors.New("databaselog driver is necessary")
	case username:
		return nil, errors.New("databaselog username is necessary")
	case password:
		return nil, errors.New("databaselog password is necessary")
	case databaseName:
		return nil, errors.New("databaselog databasename is necessary")
	}
	dbl := NewDatabaseLog(driver, username, password, ip.String(), port, databaseName)

	tables := databaselog.Tables.Table
	for _, table := range tables {
		tableName := table.Name
		if tableName != "" {
			attributes := table.Attributes.Attribute
			if len(attributes) == 0 {
				return nil, errors.New("databaselog table "+tableName+"'attribute is necessary")
			}
			cols, colErr := parseAttributes(attributes)
			if colErr != nil {
				return nil, colErr
			}
			logTable := newLogTableXML(tableName, cols)
			logTableSelfLevel, ltslErr := parseLogLevel(table.SelfLevel)
			if ltslErr != nil {
				return nil, errors.New("databaselog table "+tableName+"'s selflevel cannot be parsed")
			}
			logTable.SetTableSelfLevel(logTableSelfLevel)
			dbl.AddTable(logTable)
		}
	}
	logger.AddDatabaseLog(dbl)
	return logger, nil
}

func configSmtpLog(smtplg smtplog, logger *Logger)(*Logger, error)  {
	host := smtplg.Host
	port := smtplg.Port
	password := smtplg.Password
	from := smtplg.From
	switch "" {
	case host:
		return nil, errors.New("smtplog host is necessary")
	case password:
		return nil, errors.New("smtplog password is necessary")
	case from:
		return nil, errors.New("smtplog 'from' is necessary")
	}
	if !verifyEmailFormat(from){
		return nil, errors.New("smtplog 'from' format error")
	}
	sl := NewSmtpLog(host, password, from, port, smtplg.Subject)
	selfLevel := smtplg.SelfLevel
	if selfLevel != "" {
		smtpLogSelfLevel, smtpLevelErr := parseLogLevel(selfLevel)
		if smtpLevelErr != nil {
			return nil, errors.New("cannot parse smtolog selflevel")
		}
		sl.SetMailSelfLogLevel(smtpLogSelfLevel)
	}
	recipients := smtplg.Recipients.Recipient
	if len(recipients)<= 0 {
		return nil, errors.New("smtplog recipent is necessary")
	}
	for _, recpt := range recipients {
		if !verifyEmailFormat(recpt){
			return nil, errors.New("email recipient "+recpt+" format error")
		}
	}
	sl.SetRecipient(smtplg.Recipients.Recipient)
	logger.AddSmtpLog(sl)
	return logger, nil
}

func parseFormatRef(ref string)(*format, bool)  {
	if ref != "" {
		frmt := formatMap[ref]
		if !reflect.DeepEqual(frmt, formatXML{}) {
			if frmt.Format !="" {
				return NewFormat(frmt.Format), true
			}
		}
	}
	return nil, false
}

func checkErr(_ *Logger, err error) {
	if err != nil {
		panic(err)
	}
}