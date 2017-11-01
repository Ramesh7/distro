package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/airbrake/gobrake"
	"github.com/gin-gonic/contrib/newrelic"
	"github.com/gin-gonic/gin"
)

func init() {
	log.Formatter = new(logrus.JSONFormatter)
}

var (
	dbConfigPath      = "./database_config.json"
	configurationFlag = flag.String("configuration-path", "conf.json", "Loads configuration file")
	maxStackTraceSize = 4096
	listDatabases     = "show databases"
	log               = logrus.New()
)

type Cluster struct {
	Deployment string `json:"deployment"`
	Dsn        string `json:"dsn"`
}

type Configuration struct {
	BindAddress             string
	NewRelicLicenseKey      string
	NewRelicApplicationName string
	AirbrakeProjectID       string
	AirbrakeProjectKey      string
	Verbose                 string
}

var airbrake *gobrake.Notifier

func airbrakeRecovery(airbrake *gobrake.Notifier) gin.HandlerFunc {
	w := gin.DefaultWriter
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				rvalStr := fmt.Sprint(rval)
				w.Write([]byte(fmt.Sprintf("recovering for error:%s from uri:%s", rvalStr, c.Request.URL)))
				ginErrorHandler("Recovery", errors.New(rvalStr), c, true, true)
			}
			defer airbrake.Flush()
		}()
		c.Next()
	}
}

func main() {
	log.Info("Starting from main method...")

	log.Info("Loading DB configuration...")
	conf, err := loadConfiguration(*configurationFlag)
	if err != nil {
		log.Error(fmt.Sprintf("Error loading configuration: %v", err))
		return
	}
	verbose, err := strconv.ParseBool(conf.Verbose)
	if err != nil {
		log.Error(fmt.Sprintf("Error parsing verbose flag: %v", err))
		return
	}
	airbrakeProjectID, err := strconv.ParseInt(conf.AirbrakeProjectID, 10, 64)
	if err != nil {
		log.Error(fmt.Sprintf("Error parsing airbrake option: %v", err))
		return
	}
	airbrake = gobrake.NewNotifier(airbrakeProjectID, conf.AirbrakeProjectKey)

	for _, cluster := range clusterList {
		log.Info(cluster.Deployment)
	}

	airbrake = gobrake.NewNotifier(airbrakeProjectID, conf.AirbrakeProjectKey)

	r := gin.New()
	r.Use(airbrakeRecovery(airbrake)) // Use airbrakeRecovery as early as possible
	r.Use(newrelic.NewRelic(conf.NewRelicLicenseKey, conf.NewRelicApplicationName, verbose))
	r.Use(gin.Logger())
	buildRoutes(r)
	r.Run(conf.BindAddress)
}

func buildRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/clusters", getClusterList())
		v1.GET("/health", getHealth())
		v1.GET("/query/download", getQueryResult())
	}
}
