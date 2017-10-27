package main

import (
  "database/sql"
  "fmt"

  "github.com/gin-gonic/gin"
  _ "github.com/go-sql-driver/mysql"
)

func main() {
  fmt.Print("Starting from main method\n")
  clusterList := loadConfig()

  gin.SetMode(gin.ReleaseMode)
  router := getRouter()

  for _, cluster := range clusterList {
    fmt.Println(cluster.Deployment)
  }
  fmt.Print("End of program...")
}


func getRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// API routes
	v1 := r.Group("/v1")
	{
		v1.GET("/clusters", getClusterList)
		v1.GET("/health", getHealth)
    v1.GET("/query/download", getQueryResult)
	}
	return r
}
