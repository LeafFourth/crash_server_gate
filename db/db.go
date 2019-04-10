package db

import (
	"database/sql"
	"strconv"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"

	"crash_server_gate/common"
	"crash_server_gate/defines"
)

var conn *sql.DB;

func GetConn() *sql.DB {
	return conn;
}

func GetTableName(date time.Time) string {
	y, m, d := date.Date();
	ys := strconv.FormatInt(int64(y), 10);
	ms := strconv.FormatInt(int64(m), 10);
	ds := strconv.FormatInt(int64(d), 10);

	if len(ms) < 2 {
		ms = "0" + ms;
	}

	if len(ds) < 2 {
		ds = "0" + ds;
	}

	return "crash_" + ys + ms + ds; 
}

func PreCreateTableForDate(date time.Time) {
	if conn == nil {
		common.ErrorLogger.Print("conn nil");
		return;
	}

	name := GetTableName(date);
	sqlcmd := "CREATE table " + name + " (name VARCHAR(100) PRIMARY KEY, uid int, callstack VARCHAR(2000));"

	conn.Exec(sqlcmd);
}



func InitDB() {
	var err error;
	conn, err = sql.Open("mysql", defines.DBConnString);
	if err != nil {
	  	common.ErrorLogger.Print(err);
	  	return;
	}
}