package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	// "github.com/joho/sqltocsv"
)

var (
	clusterList = loadDBConfiguration()
)

type InputParams struct {
	InputClusterList []string
	Query            string
}

func getHealth() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

func getClusterList() func(*gin.Context) {
	return func(c *gin.Context) {
		cluster := []string{}
		for _, c := range clusterList {
			cluster = append(cluster, c.Deployment)
		}
		c.JSON(http.StatusCreated, cluster)
	}
}

func getQueryResult() func(*gin.Context) {
	return func(c *gin.Context) {
		var inputParam InputParams
		err := c.BindJSON(&inputParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		for _, cluster := range inputParam.InputClusterList {
			clusterInfo := getClusterInfo(cluster)
			if (Cluster{}) == clusterInfo {
				fmt.Println("Cluster " + cluster + " information is missing in config")
				continue
			}

			var dbConfig DbConfig
			log.Info("DSN string : " + clusterInfo.Dsn)
			dbConfig.Dsn = clusterInfo.Dsn
			dbConfig.DbName = "mysql"
			dbConfig.MaxIdleConnections = "10"
			dbConfig.MaxOpenConnections = "128"
			// getting list of databases
			_, _ = getDatabaseList(&dbConfig)
			log.Info("end of call")

			// TODO : excute query
			// 1. build DB connection
			// 2. Get all schemas
			// 3. fire query across all schemas
		}

		c.JSON(http.StatusCreated, clusterList)
	}
}

func ginErrorHandler(message string, err error, c *gin.Context, printStack bool, sendAirbrake bool) {
	w := gin.DefaultWriter
	w.Write([]byte(fmt.Sprintf("%s error:%v", message, err)))
	if printStack {
		trace := make([]byte, maxStackTraceSize)
		runtime.Stack(trace, false)
		w.Write([]byte(fmt.Sprintf("stack trace--\n%s\n--", trace)))
	}
	if sendAirbrake {
		airbrake.Notify(fmt.Errorf("%s error:%v", message, err), c.Request)
		defer airbrake.Flush()
	}
	c.AbortWithError(http.StatusInternalServerError, err)
}

func getClusterInfo(name string) Cluster {
	var cluster Cluster
	for _, c := range clusterList {
		if c.Deployment == name {
			cluster = c
		}
	}
	return cluster
}
