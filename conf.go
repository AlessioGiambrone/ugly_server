package main

import (
	"os"

	"github.com/jinzhu/configor"
)

type Constraint struct {
	Round    *int     `json:",omitempty"`
	Max, Min *float64 `json:",omitempty"`
}

var Config = struct {
	Port           int    `default:"7070"`
	ProxiedService string `default:"http://localhost:8000"`
	Constraints    map[string]Constraint
}{}

// Gets the variable from the environment. `def` is the default value
// that gets used if no env is found with that name.
func getenv(varName, def string) string {
	if newVar := os.Getenv(varName); newVar != "" {
		return newVar
	}
	return def
}

func loadConfig() {
	configor.Load(&Config, getenv("CONFIG", "config.yaml"), "conf/config.yaml")
}
