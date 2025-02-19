package server

import (
    "os"
    "bytes"

    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
    highlighting "github.com/yuin/goldmark-highlighting/v2"
    "go.abhg.dev/goldmark/mermaid"
)

func PreparePost(path string) (*bytes.Buffer, error) {
    // assumes existing '.md' file, should be checked in calling function
    raw, err := os.ReadFile(path)
    if nil != err {
	return nil, err
    }

    mdRenderer := goldmark.New(
	goldmark.WithExtensions(
	    highlighting.NewHighlighting(
		highlighting.WithStyle("dracula"),
	    ),
	    &mermaid.Extender{},
	    extension.Footnote,
	),
    )
    var parsed bytes.Buffer
    if err = mdRenderer.Convert(raw, &parsed); nil != err {
	return nil, err
    }

    return &parsed, nil

}

// TODO::
// build post from a directory with:
// <directory_name> : post name CamelCase
// - info.json : contains all relevant metadata (writer, created, updated)
// - post.md : the post itself
// - rsrc_<> : resource files (images, giffs, videos)

// TODO::
// load all posts, and update list every sync
// use channels and goroutines for updating from `sync.go`
