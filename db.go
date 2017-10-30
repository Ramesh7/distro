package main

import (
  "fmt"
  "database/sql"
  "strconv"

  _ "github.com/go-sql-driver/mysql"
)

func connection(conf *Configuration, dbName string) (db *sql.DB, err error) {
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

func executeQuery(dbName string, query string) (rows string, err error) {
  rows, err = db.Prepare(upgradeSelect)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating statement dictionarySelect: %v", err))
		return "", err
	}
	return result, err
}
