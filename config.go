package main

import (
  "flag"
  "encoding/json"
	"fmt"
	"io/ioutil"
	"os"

  "github.com/Sirupsen/logrus"
)

var (
  logger     = logrus.New()
  configPath = "./config.json" // "/etc/distro/config.json"
)

type Cluster struct {
    Deployment string `json:"deployment"`
    Slave      string `json:"slave"`
    Username   string `json:"username"`
    Password   string `json:"password"`
}

func (c Cluster) toString() string {
    return toJson(c)
}

func toJson(c interface{}) string {
    bytes, err := json.Marshal(c)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    return string(bytes)
}

func loadConfiguration() ([]Cluster, error) {
  fmt.Printf("Hello from loadConfig\n")
  flag.StringVar(&configPath, "c", configPath, "distro -c 'filePath'")
  flag.Parse()

	data, err := ioutil.ReadFile(configPath)

	if err != nil {
		logger.Errorf("Error! %s is required: %v", configPath, err)
		fmt.Printf("Error! %s is required: %v\n", configPath, err)
		os.Exit(1)
	}

  var cluster []Cluster
  err = json.Unmarshal(data, &cluster)
	if err != nil {
		logger.Errorf("Error loading %s: %v", configPath, err)
		fmt.Printf("Error loading %s: %v\n", configPath, err)
		os.Exit(1)
	}
  return cluster, nil
}
