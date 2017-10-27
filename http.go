package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getHealth(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func getClusterList(*gin.Context) {
  return func(c *gin.Context) {
		deploymentList := []string{"lena", "austin", "foo"}
    c.JSON(http.StatusCreated, deploymentList)
  }
}

func getQueryResult(c *gin.Context) {
  return func(c *gin.Context) {
    deploymentList := c.PostForm("deployments")
		query := c.PostForm("query")
    c.JSON(http.StatusOK, gin.H{"status": "success"})
  }
}
