package gologger

type LogAttrubutes string

const (
	Level LogAttrubutes = "Level"
	Time LogAttrubutes = "Time"
	Date LogAttrubutes = "Date"
	Msg LogAttrubutes = "Msg"
	Func LogAttrubutes = "Func"
	File LogAttrubutes = "File"
	Line LogAttrubutes = "Line"
)
type logTable struct {
	tableName string
	tableSelfLevel logLevel
	cols []LogAttrubutes
}

func NewLogTable(tableName string, cols ...LogAttrubutes)*logTable  {
	return &logTable{
		tableName: tableName,
		cols: cols,
		tableSelfLevel: DEBUG,
	}
}
func newLogTableXML(tableName string, cols []LogAttrubutes)*logTable  {
	return &logTable{
		tableName: tableName,
		cols: cols,
		tableSelfLevel: DEBUG,
	}
}

func (lt *logTable)SetTableSelfLevel(level logLevel)  {
	lt.tableSelfLevel = level
}
