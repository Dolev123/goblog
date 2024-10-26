package config

import (
	"encoding/json"
	"log"
	"os"

	"fmt"
)

type Config struct {
	ListenAddr  string `json:"address"`
	Method      string `json:"method"`
	Source      string `json:"source"`
	Destination string `json:"dest"`
	Secrets     string `json:"secrets"`
}

func LoadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if nil != err {
		log.Fatal("could not extract configuration from `"+path+"`: ", err)
	}
	var config *Config
	if nil != json.Unmarshal(data, &config) {
		log.Fatal("could not unmarshal configuration from `"+path+"`: ", err)
	}
	return config
}

func DebugConfig(conf *Config) {
	fmt.Println("ListenAddr:", conf.ListenAddr)
	fmt.Println("Method:", conf.Method)
	fmt.Println("Source:", conf.Source)
	fmt.Println("Destination:", conf.Destination)
	fmt.Println("Secrets:", conf.Secrets)
}
