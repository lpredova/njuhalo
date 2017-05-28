package configuration

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lpredova/shnjuskhalo/model"
)

const configFile = ".njhalo.json"

// CreateConfigFile creates empty configuration file in cwd
func CreateConfigFile() bool {

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
