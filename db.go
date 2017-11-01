package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	Dsn                string
	DbName             string
	MaxOpenConnections string
	MaxIdleConnections string
}

func openDbConnection(conf *DbConfig, dbName string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", conf.Dsn+dbName)
	if err != nil {
		log.Error(fmt.Sprintf("Error opening mysql connection: %v", err))
		return nil, err
	}
	maxIdleConnections, err := strconv.Atoi(conf.MaxIdleConnections)
	if err != nil {
		log.Error("Invalid entry for MaxIdleConnections")
		return nil, err
	}
	db.SetMaxIdleConns(maxIdleConnections)
	maxOpenConnections, err := strconv.Atoi(conf.MaxOpenConnections)
	if err != nil {
		log.Error("Invalid entry for MaxOpenConnections")
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConnections)
	return db, nil
}

func getDatabaseList(conf DbConfig, cluster string) (dbList []string, err error) {
	db, err := openDbConnection(&conf, "mysql")
	if err != nil {
		log.Error("Unable to make connection for " + cluster + "cluster")
	}

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		log.Error(fmt.Sprintf("Error listing database: %v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var Database string
		rows.Scan(&Database)
		matched, _ := regexp.MatchString("^mysql$|^test$|^information_schema$|^performance_schema$|^sys$|^sftp$|^klosetd$|^resque_shepherd$|^shasta$|^percona$", Database)
		if matched {
			continue
		}
		dbList = append(dbList, Database)
	}
	return dbList, nil
}

func getQueryResult(conf DbConfig, query string, dbName string, cluster string) (result [][]string, err error) {
	db, err := openDbConnection(&conf, dbName)
	if err != nil {
		log.Error("Unable to make connection for " + cluster + "cluster")
		return result, err
	}

	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return result, err
	}

	log.Info(reflect.TypeOf(rows))
	columnNames, err := rows.Columns()
	if err != nil {
		return result, err
	}

	log.Info(reflect.TypeOf(columnNames))
	//columns, err := rows.Columns()
	//readCols := make([]interface{}, len(columns))
	log.Info("============")
	//og.Info(reflect.TypeOf(readCols))
	return result, err
}
