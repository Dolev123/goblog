package server

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/ast"
    "github.com/yuin/goldmark/renderer"
    "github.com/yuin/goldmark/renderer/html"
    "github.com/yuin/goldmark/util"
)

const (
    RemoteIndicator = "://"
)

func NewPostImageOption(postId int) goldmark.Option {
    return goldmark.WithRendererOptions(
	renderer.WithNodeRenderers(
	    util.Prioritized(NewPostImageRenderer(postId), 0),
	),
    )
}

type postImageRenderer struct {
    html.Config
    postId int
}

func NewPostImageRenderer(postId int, options ...html.Option) goldmark.Extender {
    config := html.NewConfig()
    for _, opt := range options {
	opt.SetHTMLOption(&config)
    }
    return &postImageRenderer{
	Config: config,
	postId: postId,
    }
}

func (r *postImageRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
    reg.Register(ast.KindImage, r.renderImage)
}

// add lazy loading and post resource hirarchy (if it exists).
func (r *postImageRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
    if !entering {
	return ast.WalkContinue, nil
    }
    
    n := node.(*ast.Image)
    w.WriteString("<img src=\"")
    if r.Unsafe || !html.IsDangerousURL(n.Destination) {
	// add postId only if it is a valid value, and not a remote link
	if isRemoteResource(string(n.Destination)) {
	    n.SetAttribute([]byte("title"), fmt.Sprintf("Source: %s", n.Destination))
	} else if r.postId >= 0 && len(postsMetadata) > r.postId {
	    w.WriteString(strconv.Itoa(r.postId) + "/")
	}
	w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
    }
    w.WriteString(fmt.Sprintf(
	"\" alt=\"%s\" loading=\"lazy\" class=\"post-image\"",
	n.Text(source),
    ))
    if n.Attributes() != nil {
	html.RenderAttributes(w, n, html.ImageAttributeFilter)
    }
    w.WriteString("/>")
    return ast.WalkSkipChildren, nil
}

// Implement goldmark.Extender interface
func (r *postImageRenderer) Extend(m goldmark.Markdown) {
    m.Renderer().AddOptions(
	renderer.WithNodeRenderers(
	    util.Prioritized(r, 0),
	),
    )
}

// Check if a resource is from a remote location.
func isRemoteResource(resource string) bool {
    return strings.Contains(resource, RemoteIndicator)
}
