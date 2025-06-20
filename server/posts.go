package server

import (
    "bytes"
    "os"

    "github.com/yuin/goldmark"
    highlighting "github.com/yuin/goldmark-highlighting/v2"
    "github.com/yuin/goldmark/extension"
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
