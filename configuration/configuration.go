package configuration

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lpredova/njuhalo/model"
)

const interval = 5
const sleep = 3
const configFile = "config.json"

const path = "./storage/" + configFile

// ParseConfig parsers currently available config file
func ParseConfig() model.Configuration {
	var configuration = model.Configuration{}
	file, err := loadFileConfig()
	if err != nil {
		return configuration
	}

	decoder := json.NewDecoder(file)
	decoder.Decode(&configuration)
	return configuration
}

// PrintConfig prints configuration file
func PrintConfig() {
	var configuration = ParseConfig()
	configJSON, err := json.MarshalIndent(configuration, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(configJSON))
}

// CreateFileConfig creates empty configuration file in cwd
func CreateFileConfig(conf model.Configuration) bool {

	conf.RunIntervalMin = interval
	conf.SleepIntervalSec = sleep

	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer f.Close()

	jsonConfig, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	f.WriteString(string(jsonConfig))
	return true
}

// Load config json file, in cwd and then user home folder
func loadFileConfig() (*os.File, error) {

	file, err := os.Open(configFile)
	if err != nil {
		file, err = os.Open(path)
		return file, err
	}

	return file, nil
}
