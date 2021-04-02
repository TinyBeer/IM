package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var sqldb *sql.DB

func Init(username string, password string, addr string, port int, dbName string) (err error) {
	// connect to mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		username,
		password,
		addr,
		port,
		dbName)
	// dataSourceName
	sqldb, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("failed to open database, err:" + err.Error())
	}

	err = sqldb.Ping()
	if err != nil {
		panic("failed to connect to mysql database, err:" + err.Error())
	}
	return nil
}

func GetSqlDB() *sql.DB {
	return sqldb
}
