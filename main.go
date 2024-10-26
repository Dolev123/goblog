package main

import (
	"flag"

	"github.com/Dolev123/goblog/config"
	"github.com/Dolev123/goblog/router"
	"github.com/Dolev123/goblog/sync"
)

func main() {
	fconf := flag.String("config", "config.json", "Path to JSON configuration file")
	flag.Parse()

	conf := config.LoadConfig(*fconf)
	config.DebugConfig(conf)
	sync.SyncPosts(conf)
	router.StartServer(conf)
}
