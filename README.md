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
    "address":   "[ip]:[port]",
    "method":    "git|directory",
    "source":    "git://repo|/src/path/",
    "dest":      "/path/to/server/files",
    "schedule":  "@cron",
    "secrets":   "/path/to/secrets.json",
    "structure": "bare|full",
    "title":     "Example Blog Title"
}
```

- The `dest` directory will be the directory for the server.  
`source` can be either a repo address or a full path to a directory.  
- `schedule` is based on cron schedule. for more info check out [here](https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format)  
- Server running with `full` structure expects to have a directory for each post, and a `"resources"` directory. And each post's directory should conatin atleast `"metadata.json"` and a markdown file for the post's content. The `bare` structure expects that all the posts are contained inside the same directory, and allows for only markdown files. `bare` is recommended only for testing.

The `metadata.json` file has the follwing structure:
```json
{
    "author": "<your name>",
    "created": "yyyy-MM-dd hh:mm:ss",
    "updated": "yyyy-MM-dd hh:mm:ss",
    "title": "<Post Title>"
}
```

## Running with Dockerfile
The `Dockerfile` expects to have a `config.json` when building the docker image.   

For running the examles located inside the `example` directory you should do the following changes:  
1. Change directory to the desired example (e.g. `cd example/full_example`).
2. Copy the Dockerfile to the current directory `cp ../../Dockerfile .`. 
3. Edit the local Dockerfile, and add the following to the "setup the site" section: `ADD . /src/path`.
4. Build and Run the docker.

## Creating HTML Templates

This section applies only to full mode/structure.
The html is set up using golang's builtin `html/template`. If you want to change them, you can, and have some parameters to access:
- `{{.ID}}` - calculated post ID
- `{{.Created}}` and `{{.Updated}}` - dates from metadata
- `{{.Author}}` - author from metadata
- `{{.Title}}` - title from metadata
- `{{.Content}}` - the html generated from the markdown.
- `{{.BlogTitle}}` - title from metadata

The values may be accessed a bit differently, based on file's location, and expects exactly the files listed under resources.

## Formating the code

Since golang has it's own builtin standart formatiing tools, I do use it, but change the tabs to 4 spaces, because I like it more.  
For thise interested, here is the command:
```sh
find . -regex '.*.go$' -exec go fmt {} \; -exec sed -i -e 's/\t/    /' {} \;
# for each go file, run 'go fmt' and then 'sed'.
```
