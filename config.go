package main

import (
  "fmt"
  "encoding/json"
	"io/ioutil"
	"os"
)

func (c Cluster) toString() string {
    return toJson(c)
}

func toJson(c interface{}) string {
    bytes, err := json.Marshal(c)
    if err != nil {
        log.Error(err.Error())
        os.Exit(1)
    }

    return string(bytes)
}

func loadDBConfiguration() ([]Cluster) {

	data, err := ioutil.ReadFile(dbConfigPath)

	if err != nil {
		log.Error("Error! %s is required: %v", dbConfigPath, err)
		os.Exit(1)
	}

  var cluster []Cluster
  err = json.Unmarshal(data, &cluster)
	if err != nil {
		log.Error("Error loading %s: %v", dbConfigPath, err)
		os.Exit(1)
	}
  return cluster
}

func loadConfiguration(configPath string) (*Configuration, error) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Error(fmt.Sprintf("Error opening config file:%s error:%v", configPath, err))
		return nil, err
	}
	decoder := json.NewDecoder(file)
	var configuration Configuration
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Error(fmt.Sprintf("Error decoding config file:%s error:%v", configPath, err))
		return nil, err
	}
	return &configuration, nil
}
