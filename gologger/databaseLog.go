package gologger

import (
	"database/sql"
	"fmt"
)

type databaseLog struct {
	databaseName string
	db *sql.DB
	logTables []*logTable
}

func NewDatabaseLog(driver, username, passwd, ip string, port int64, databaseName string)*databaseLog  {
    db , openErr:= sql.Open(driver, username+":"+passwd+"@tcp("+ip+":"+fmt.Sprintf("%d", port)+")/"+databaseName+"?charset=utf8")
    if openErr != nil {
    	panic(openErr)
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	//db.SetConnMaxLifetime()
	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}
	return &databaseLog{
		databaseName: databaseName,
		db: db,
	}
}

func (dbl *databaseLog)AddTable(table *logTable)  {
	judgeSql := "SELECT table_name FROM information_schema.TABLES WHERE table_name ='"+table.tableName+"';"
	rows, _ := dbl.db.Query(judgeSql)
	t, _ :=rows.Columns()
	if len(t) == 0 {
		idcol := "id int primary key auto_increment not null,"
		var otherCols string
		for _, col := range table.cols{
			otherCols = otherCols + fmt.Sprintf(" %s", col)+" varchar(255),"
		}
		cols := idcol+otherCols
		cols = cols[:len(cols)-1]
		_, createErr := dbl.db.Exec("create table "+table.tableName+"("+cols+")")
		if createErr != nil {
			panic(createErr)
		}
	}
	dbl.logTables = append(dbl.logTables, table)
}

func (dbl *databaseLog)insertLog(msg *logMsg)  {
	for _, table := range dbl.logTables {
		var colVals string

		msgLevel, logErr := parseLogLevel(msg.levStr)
		if logErr != nil {
			panic(logErr)
		}
		if table.tableSelfLevel <= msgLevel {
			for _, col := range table.cols {
				colVals = colVals+",'"+ getParsedValue(fmt.Sprintf("%s", col), msg)+"'"
			}
			insertSql := "insert into "+table.tableName +" values (default"+colVals+")"
			_, insertEerr := dbl.db.Exec(insertSql)
			if insertEerr != nil {
				panic(insertEerr)
			}
		}
	}
}