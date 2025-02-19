package sync

import (
    "os"
    "os/exec"

    "github.com/robfig/cron"

    "github.com/Dolev123/goblog/config"
    pkglog "github.com/Dolev123/goblog/logger"
)

/*
NOTE assumes the following structure of files:
<server-base>
    - blog1.md
    - blog2.md
    ...
In the future, probably it is prefered to be:
<server-base>
    - blog1_dir
	- README.md
	- blog1.resource1
	- blog1.resource2
	...
    - blog2_dir
    ...
*/

var logger = pkglog.CreateNewLogger()

func SyncPosts(conf *config.Config, syncChans []chan bool) {
    // check for destination
    if _, err := os.ReadDir(conf.Destination); nil != err {
	if os.IsNotExist(err) {
	    logger.Println("Destination directory `" + conf.Destination + "` does not exist, creating it")
	    if nil != os.MkdirAll(conf.Destination, os.ModePerm) {
		logger.Fatal("Failed creating destination directory")
	    }
	}
    }
    // sync
    switch conf.Method {
    case "directory":
	logger.Println("Calling directory syncronization")
	directorySync(conf)
    case "git":
	logger.Println("Calling git syncronization")
	gitSync(conf)
    default:
	logger.Fatal("Unknown method for syncronization. Aborting...")
    }
    // update all related sync channels
    logger.Println("Sending update signals")
    for i, ch := range syncChans {
	logger.Println("sending to channel indexed:", i)
	go func(){ch <- true}()
    }
}

func StartCronSync(conf *config.Config, syncChans []chan bool) *cron.Cron {
    cr := cron.New()
    cr.AddFunc(conf.Schedule, func(){SyncPosts(conf, syncChans)})
    cr.Start()
    return cr
}

func directorySync(conf *config.Config) {
    // TODO:: check for trailing "/"
    src := conf.Source
    dst := conf.Destination

    entries, err := os.ReadDir(src)
    if nil != err {
	if os.IsNotExist(err) {
	    logger.Println("Source directory does not exist:", src)
	    return
	}
	logger.Println("Unknown error while reading src directory:", err)
	return
    }

    for _, entry := range entries {
	fname := entry.Name()
	cmd := exec.Command("cp", "-r", src+fname, dst+fname)
	cmd.Run()
	// TODO:: add log / goncurrency?
    }
}

func gitSync(conf *config.Config) {
    repo := conf.Source
    dst := conf.Destination

    should_clone := false

    // check if already cloned repository
    cmd := exec.Command("git", "-C", dst, "rev-parse", "--is-inside-work-tree")
    if err := cmd.Run(); nil != err {
	logger.Println(":DEBUG:", err)
	if exiterr, ok := err.(*exec.ExitError); ok && 0 != exiterr.ExitCode() {
	    should_clone = true
	} else {
	    logger.Println("Failed to determine if repository exists")
	    return
	}
    }

    if should_clone {
	logger.Println("Cloning git repo")
	cmd = exec.Command("git", "-C", dst, "clone", repo,  ".")
    } else {
	logger.Println("Updateing (pull) git repo")
	cmd = exec.Command("git", "-C", dst, "pull", "origin")
    }
    if err := cmd.Run(); nil != err {
	logger.Println("Failed to run git command:", err) 
    }
}
