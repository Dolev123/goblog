package test

import (
    "path/filepath"
    "strings"
    "testing"
    "time"

    "github.com/Dolev123/goblog/server"
    "github.com/Dolev123/goblog/types"
)

func TestExamplePostMetadata(t *testing.T) {
    exampleMetadataDirPath := filepath.Join("..", "example")
    metadata, err := server.LoadPostMetada(exampleMetadataDirPath)
    if err != nil {
    	t.Fatal(err.Error())
    }
    t.Logf("Metadata: %v", metadata)
    if metadata.Writer != "ME" {
    	t.Error("writer", metadata.Writer)
    }
    want_created := types.Time{time.Date(2024, 12, 17, 01, 14, 00, 00, time.UTC)}
    if metadata.Created != want_created {
    	t.Error("created", metadata.Created)
    }
    want_updated := types.Time{time.Date(2025, 01, 18, 02, 15, 01, 00, time.UTC)}
    if metadata.Updated != want_updated {
    	t.Error("updated", metadata.Updated)
    }
    if metadata.Title != "Test Post" {
    	t.Error("title", metadata.Title)
    }
    if metadata.Path != exampleMetadataDirPath {
    	t.Error("path", metadata.Path)
    }

}

func TestConvertTitleToPath(t *testing.T) {
    const correct = "test_title.md"
    tests := []string{
    	"test_title", "Test Title", " Test Title", "\tTest\tTitle\t", "TEST TITLE",
    }
    for _, test := range tests {
    	path := server.ConvertTitleToPath(test)
    	if correct != path {
    		t.Error("got:", path, "instead of:", correct)
    	}
    }
}

func TestLoadAndRenderPost(t *testing.T) {
    const correct = `<h1>Test :)</h1>
<h2>This is paragraph 1</h2>
<p>This is a sentence.<br>
<em><strong>this is italic bold</strong></em></p>
<ul>
<li>[ ] item 1</li>
<li>[X] item 2</li>
</ul>
<h2>This is paragraph 2</h2>
<pre tabindex="0" style="color:#f8f8f2;background-color:#282a36;"><code><span style="display:flex;"><span><span style="color:#8be9fd;font-style:italic">echo</span> <span style="color:#f1fa8c">&#39;12&#39;</span>
</span></span></code></pre><pre tabindex="0" style="color:#f8f8f2;background-color:#282a36;"><code><span style="display:flex;"><span><span style="color:#ff79c6">def</span> <span style="color:#50fa7b">foo</span>():
</span></span><span style="display:flex;"><span>    <span style="color:#8be9fd;font-style:italic">print</span>(<span style="color:#f1fa8c">&#34;?&#34;</span>, <span style="color:#bd93f9">12</span>)
</span></span></code></pre>`
    t.Logf("Correct: %v", correct)
    exampleMetadataDirPath := filepath.Join("..", "example")
    metadata, err := server.LoadPostMetada(exampleMetadataDirPath)
    if err != nil {
    	t.Fatal("Failed to load metadata:", err.Error())
    }
    post, err := server.LoadAndRenderPostData(metadata)
    if err != nil {
    	t.Fatal("Failed to load/render post:", err.Error())
    }
    data := post.Data().String()
    if correct != data {
    	t.Logf("Compare(data, correct): %v", strings.Compare(data, correct))
    	t.Error("Incorrect Data:", data)
    }
}
