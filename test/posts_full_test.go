package test

import (
    "testing"
    "path/filepath"
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
    if metadata.Created !=  want_created {
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
