package main

import (
  "fmt"

  "github.com/gin-gonic/gin"
  "github.com/Sirupsen/logrus"
)

var (
  log = logrus.New()
)

type Configuration struct {
	Dsn                     string
	DbName                  string
	BindAddress             string
	MaxIdleConnections      string
	MaxOpenConnections      string
}

func main() {
  fmt.Print("Starting from main method\n")
  clusterList, err := loadConfiguration()

  if err != nil {
		log.Error(fmt.Sprintf("Error loading configuration: %v", err))
		return
	}

  for _, cluster := range clusterList {
    fmt.Println(cluster.Deployment)
  }

  r := gin.New()
  r.Use(gin.Logger())
	buildRoutes(r)
	r.Run("8080")

  fmt.Print("End of program...")
}

func buildRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
    v1.GET("/clusters", getClusterList())
    v1.GET("/health", getHealth())
    v1.GET("/query/download", getQueryResult())
	}
}
