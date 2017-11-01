package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/joho/sqltocsv"
)

var (
	clusterList = loadDBConfiguration()
	DBCon       *sql.DB
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

func getResult() func(*gin.Context) {
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
			dbConfig.Dsn = clusterInfo.Dsn
			dbConfig.MaxIdleConnections = "10"
			dbConfig.MaxOpenConnections = "128"

			dbList, err := getDatabaseList(dbConfig, cluster)

			if err != nil {
				log.Error("Unable to get database list for " + cluster + "cluster")
				continue
			}
			log.Info(strings.Join(dbList, ","))
			for _, d := range dbList {
				_, err := getQueryResult(dbConfig, inputParam.Query, d, cluster)
				if err != nil {
					log.Error("Unable to get result for " + d + " database")
					log.Error(err)
					continue
				}
			}
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
