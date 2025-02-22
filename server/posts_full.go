package server

import (
    "encoding/json"
    "os"
    "bytes"
    "path/filepath"
    "strings"
    "html/template"
    "net/http"
    "errors"
    "fmt"

    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
    highlighting "github.com/yuin/goldmark-highlighting/v2"
    "go.abhg.dev/goldmark/mermaid"

    "github.com/Dolev123/goblog/types"
)

const (

)

var (
    postTmpl, indexTmpl *template.Template
    
    postsMetadata = []*PostMetadata{}
)

func LoadTemplates(base string) error {
    base = filepath.Join(gconf.Destination, base)
    tmpPostTmpl, error := template.ParseFiles(
	filepath.Join(base, "post.html.tpl"),
	filepath.Join(base, "header.html.tpl"),
	filepath.Join(base, "footer.html.tpl"),
    )
    if nil != error {
	logger.Println("Failed to load Post template:", error)
	return error
    }
    tmpIndexTmpl, error := template.ParseFiles(
	filepath.Join(base, "index.html.tpl"),
	filepath.Join(base, "preview.html.tpl"),
	filepath.Join(base, "header.html.tpl"),
	filepath.Join(base, "footer.html.tpl"),
    )
    if nil != error {
	logger.Println("Failed to load Post template:", error)
	return error
    }
    postTmpl, indexTmpl = tmpPostTmpl, tmpIndexTmpl
    return nil
}

type Post struct {
    metadata *PostMetadata
    data *bytes.Buffer
}

func (p *Post) Data() *bytes.Buffer {
    return p.data
}

type PostMetadata struct {
    Author string `json:"writer"`
    Created types.Time `json:"created"`
    Updated types.Time `json:"updated"`
    Title string `json:"title"`
    // fields not in json
    Path string `json:"-"`
    id int `json:"-"` // to be set by loading function
}

func (metadata * PostMetadata) setID(id int) {
    metadata.id = id
}

func (metadata * PostMetadata) ID() int {
    return metadata.id 
}

func LoadPostMetada(base string) (*PostMetadata, error) {
    base = filepath.Join(gconf.Destination, base)
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

func LoadAllMetadata() error {
    var tmpMetadata []*PostMetadata
    entries, err := os.ReadDir(gconf.Destination)
    if nil != err {
	logger.Println("Failed reading directory ", entries)
	return err
    }
    currentId := 0
    for _, entry := range entries {
	if !entry.IsDir() || entry.Name() == "resources"{
	    continue
	}
	metadata, err := LoadPostMetada(entry.Name())
	if err != nil {
	    logger.Println("Failed to load post's metadata:", entry.Name())
	    continue
	}
	logger.Println("Loaded post's metadata:", metadata.Title)
	metadata.setID(currentId)
	currentId += 1
	tmpMetadata = append(tmpMetadata, metadata)
    }
    postsMetadata = tmpMetadata
    return nil
}

func ConvertTitleToPath(title string) string {
    s := title
    s = strings.ToLower(s)
    s = strings.TrimSpace(s)
    s = strings.ReplaceAll(s, " ", "_")
    s = strings.ReplaceAll(s, "\t", "_")
    s = s + ".md"
    return s
}

func LoadAndRenderPostData(metadata *PostMetadata) (*Post, error) {
    // metadata have already been joined with gconf.Destination
    path := filepath.Join(metadata.Path, ConvertTitleToPath(metadata.Title))
    raw, err := os.ReadFile(path)
    if nil != err {
	return nil, err
    }
    
    mdRenderer := goldmark.New(
	NewPostImageOption(metadata.ID()),
	goldmark.WithExtensions(
	    highlighting.NewHighlighting(
		highlighting.WithStyle("dracula"),
	    ),
	    &mermaid.Extender{},
	    extension.Footnote,
	    extension.TaskList,
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

func WritePreviewsToResponse(w http.ResponseWriter) error {
    // TODO:: split execution and writing to request
    err := indexTmpl.Execute(w, map[string]interface{}{
	"postsMetadata": postsMetadata,
	"BlogTitle": gconf.BlogTitle,
    })
    if err != nil {
	logger.Println("Failed to execute 'indexTmpl':", err)
    }
    return err
}


func WritePostToResponse(postID int, w http.ResponseWriter) error {
    if len(postsMetadata) <=  postID {
	// TODO:: replace with custom error
	return errors.New(fmt.Sprintf("No Post Available for ID: %v", postID))
    }
    post, err := LoadAndRenderPostData(postsMetadata[postID])
    if err != nil {
	return err
    }
    err = postTmpl.Execute(w, map[string]interface{}{
	"BlogTitle": gconf.BlogTitle,
	"metadata": post.metadata,
	"Content": template.HTML(post.data.String()),
    })
    if err != nil {
	logger.Println("Failed to execute 'postTmpl':", err)
    }
    return err
}
