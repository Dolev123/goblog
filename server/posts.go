package server

import (
    "os"
    "bytes"

    "github.com/yuin/goldmark"
    highlighting "github.com/yuin/goldmark-highlighting/v2"
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
	),
    )
    var parsed bytes.Buffer
    if err = mdRenderer.Convert(raw, &parsed); nil != err {
	return nil, err
    }

    return &parsed, nil

}
