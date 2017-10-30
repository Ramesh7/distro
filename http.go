package main

import (
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
    // clustertList := c.PostForm("deployments")
		// query := c.PostForm("query")
    c.JSON(http.StatusOK, gin.H{"status": "success"})
  }
}
