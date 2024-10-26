# Simple Go Blog

A very simple (and currently ugly) blog server written in go.  
All posts are made out of Markdown files, currently without any html beuatification.  
The server pulls all posts from a "source" location which can be one of:  
- "directory": copy from one directory to the other.
- "git": copy from a git repository (prefered).

## Installation

There is a Dockerfile provided, this is the prefered way for running the server.  

If you decide to not use:  
Notice that the server works for __go 1.22__, prior versions are not supported, make sure you have the correct go version.  
```sh
# cehck version
go version
# install 
go install github.com/Dolev123/goblog@latest
```

## Configure

Configuring is done with a JSON configuration file:
```json
{
    "address":  "[ip]:[port]",
    "method":   "git|directory",
    "source":   "git://repo|/src/path/",
    "dest":     "/path/to/server/files",
    "schedule": "@cron",
    "secrets":  "/path/to/secrets.json"
}
```

The `dest` directory will be the directory for the server.  
`source` can be either a repo address or a full path to a directory.  
The `Dockerfile` expects to have a `config.json` when building the docker image.  
`schedule` is based on cron schedule. for more info check out [here](https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format)  

