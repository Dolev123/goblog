package config

import (
    "encoding/json"
    "os"

    pkglog "github.com/Dolev123/goblog/logger"
)

var logger = pkglog.CreateNewLogger()

type Config struct {
    // 'IP:Port'
    ListenAddr string `json:"address"`
    // either 'git' or 'directory'
    Method      string `json:"method"`
    Source      string `json:"source"`
    Destination string `json:"dest"`
    // cron syntax
    Schedule string `json:"schedule"`
    // '/path/to/file'
    Secrets string `json:"secrets"`
    // either 'bare' or 'full'
    Structure string `json:"structure,omitempty"`
    BlogTitle string `json:"title"`
}

func LoadConfig(path string) *Config {
    data, err := os.ReadFile(path)
    if nil != err {
    	logger.Fatal("could not extract configuration from `"+path+"`: ", err)
    }
    var config *Config
    if nil != json.Unmarshal(data, &config) {
    	logger.Fatal("could not unmarshal configuration from `"+path+"`: ", err)
    }
    return config
}

func DebugConfig(conf *Config) {
    logger.Println("ListenAddr:", conf.ListenAddr)
    logger.Println("Method:", conf.Method)
    logger.Println("Source:", conf.Source)
    logger.Println("Destination:", conf.Destination)
    logger.Println("Schedule:", conf.Schedule)
    logger.Println("Secrets:", conf.Secrets)
    logger.Println("Structure:", conf.Structure)
    logger.Println("BlogTitle:", conf.BlogTitle)
}
