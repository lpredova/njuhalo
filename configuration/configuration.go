package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"

	"github.com/lpredova/shnjuskhalo/model"
)

const configFile = "config.json"

var usr, _ = user.Current()
var path = usr.HomeDir + "/.njuhalo/" + configFile

// ParseConfig parsers currently avaliable config file
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

// CreateFileConfig creates empty configuration file in cwd
func CreateFileConfig(conf model.Configuration) bool {

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

// AppendFilterToConfig appends new filter to queries
func AppendFilterToConfig(filter model.Query) bool {
	configuration := ParseConfig()
	filters := []model.Query{}

	fmt.Println(configuration.Queries)

	for k, v := range configuration.Queries {
		fmt.Println(k)

		filters = append(filters, v)
	}
	filters = append(filters, filter)

	fmt.Println(filters)

	//configuration.Queries = append(configuration.Queries, filter)

	configuration.Queries = filters

	if CreateFileConfig(configuration) {
		return true
	}

	return false
}
