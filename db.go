package main

import (
	"database/sql"
	"fmt"
	"reflect"
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

func getDatabaseList(conf *DbConfig) (dbList []string, err error) {
	db, err := openDbConnection(conf, "mysql")
	if err != nil {
		log.Error(fmt.Sprintf("Error opening database: %v", err))
		return nil, err
	}
	rows, err := db.Query("show databases" + conf.DbName)
	if err != nil {
		log.Error(fmt.Sprintf("Error listing database: %v", err))
		return dbList, nil
	}
	columns, _ := rows.Columns()
	log.Info("---------")
	log.Info(reflect.TypeOf(columns))
	log.Info("---------")
	// log.Info("Return of all db list  : " + lists)
	return dbList, err
}
