package main

import (
	"fmt"
	"runtime"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getHealth() func(*gin.Context) {
	return func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "success"})
  }
}

func getClusterList() func(*gin.Context) {
	return func(c *gin.Context) {
		clustertList := []string{"lena", "austin", "foo"}
    c.JSON(http.StatusCreated, clustertList)
  }
}

func getQueryResult() func(*gin.Context) {
  return func(c *gin.Context) {
		clustertList := []string{"lena", "austin", "foo"}
    c.JSON(http.StatusCreated, clustertList)
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
