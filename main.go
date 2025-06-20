package main

import (
    "flag"

    "github.com/Dolev123/goblog/config"
    pkglog "github.com/Dolev123/goblog/logger"
    "github.com/Dolev123/goblog/server"
    "github.com/Dolev123/goblog/sync"
)

var logger = pkglog.CreateNewLogger()

func main() {
    fconf := flag.String("config", "config.json", "Path to JSON configuration file")
    flag.Parse()

    serverSyncChan := make(chan bool)
    var syncChannles []chan bool
    syncChannles = append(syncChannles, serverSyncChan)

    conf := config.LoadConfig(*fconf)
    config.DebugConfig(conf)
    sync.SyncPosts(conf, nil)
    sync.StartCronSync(conf, syncChannles)
    server.StartServer(conf, serverSyncChan)
}
