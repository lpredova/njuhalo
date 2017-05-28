package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"

	"github.com/lpredova/shnjuskhalo/model"
)

const configFile = ".njhalo.json"

// ParseConfig is method that parsers currently avaliable config file
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

// CreateConfigFile creates empty configuration file in cwd
func CreateFileConfig() bool {

	_, err := os.Stat(configFile)
	if err == nil {
		return false
	}

	f, err := os.Create(configFile)
	if err != nil {
		return false
	}
	defer f.Close()

	conf := &model.Configuration{}
	jsonConfig, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return false
	}

	fmt.Println(string(jsonConfig))
	f.WriteString(string(jsonConfig))
	return true
}

// Try to load config json file, in cwd and then user home folder
func loadFileConfig() (*os.File, error) {

	file, err := os.Open(configFile)
	if err != nil {

		usr, err := user.Current()
		if err != nil {
			return file, err
		}
		file, err = os.Open(usr.HomeDir + "/" + configFile)
		return file, err
	}

	return file, nil
}
