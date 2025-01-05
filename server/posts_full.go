package server

import (
    "encoding/json"
    "os"
    "bytes"
    "path/filepath"
    "strings"

    "github.com/yuin/goldmark"
    highlighting "github.com/yuin/goldmark-highlighting/v2"

    "github.com/Dolev123/goblog/types"
)

type Post struct {
    metadata *PostMetadata
    data *bytes.Buffer
}

type PostMetadata struct {
    Writer string `json:"writer"`
    Created types.Time `json:"created"`
    Updated types.Time `json:"updated"`
    // maybe:
    Title string `json:"title"`
    Path string `json:"-"`
}

func LoadPostMetada(base string) (*PostMetadata, error) {
    path := filepath.Join(base, "metadata.json")
    raw, err := os.ReadFile(path)
    if nil != err {
	return nil, err
    }

    var metadata PostMetadata
    if err = json.Unmarshal(raw, &metadata); nil != err {
	return nil, err
    }
    metadata.Path = base
    return &metadata, nil
}

func ConvertTitleToPath(title string) string {
    s := title
    s = strings.ToLower(s)
    s = strings.TrimSpace(s)
    s = strings.ReplaceAll(s, " ", "_")
    s = strings.ReplaceAll(s, "\t", "_")
    s = s + ".md"
    logger.Printf("[DEBUG] title: \"%v\" => path: \"%v\"\n", title, s)
    return s
}

func LoadAndRenderPost(metadata *PostMetadata) (*Post, error) {
    raw, err := os.ReadFile(metadata.Path)
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

    return &Post{
	metadata: metadata,
	data: &parsed,
    }, nil
}
