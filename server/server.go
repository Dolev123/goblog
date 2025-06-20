package server

import (
    "bytes"
    "fmt"
    "html"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"

    "github.com/Dolev123/goblog/config"
    pkglog "github.com/Dolev123/goblog/logger"
)

var gconf *config.Config

var logger = pkglog.CreateNewLogger()

func StartServer(conf *config.Config, syncChan chan bool) {
    gconf = conf
    mux := http.NewServeMux()
    switch gconf.Structure {
    case "bare":
    	mux.HandleFunc("GET /", handleIndexBare)
    	mux.HandleFunc("GET /{post}", handlePostBare)
    case "full":
    	setupFullMode(mux, syncChan)
    default:
    	logger.Fatal("Unknown server mode/structure:", gconf.Structure)
    }
    logger.Println("setting server in", gconf.Structure, "mode")
    srv := &http.Server{
    	Addr:    conf.ListenAddr,
    	Handler: mux,
    	// should make a new logger for error?
    	ErrorLog: logger,
    	// TODO:: add TLS
    }
    logger.Fatal("Server Failed with:", srv.ListenAndServe())
}

func setupFullMode(mux *http.ServeMux, syncChan chan bool) {
    if LoadTemplates("resources") != nil {
    	logger.Fatal("Failed initializing templates for 'full' mode. Aborting...")
    }
    if LoadAllMetadata() != nil {
    	logger.Fatal("Failed loading post's metedata for 'full' mode. Aborting...")
    }
    // start reloading corotuine
    go func() {
    	for <-syncChan {
    		if LoadTemplates("resources") != nil {
    			logger.Println("Failed reloading templates for 'full' mode. Aborting...")
    		}
    		if LoadAllMetadata() != nil {
    			logger.Println("Failed reloading post's metedata for 'full' mode. Aborting...")
    		}
    	}
    }()

    mux.HandleFunc("GET /", handleIndexFull)
    mux.HandleFunc("GET /{post}", handlePostFull)
    mux.HandleFunc("GET /{post}/{rsrc}", handlePostResourceFull)
    mux.HandleFunc("Get /resources/{rsrc}", handleResourcesFull)
}

// GET "/{post}"
func handlePostBare(w http.ResponseWriter, req *http.Request) {
    fname := req.PathValue("post")
    if !strings.HasSuffix(fname, ".md") {
    	http.Redirect(w, req, "/", http.StatusFound)
    	return
    }
    fname = filepath.Join(gconf.Destination, fname)
    finfo, err := os.Stat(fname)
    if nil != err || !finfo.Mode().IsRegular() {
    	w.WriteHeader(http.StatusNotFound)
    	fmt.Fprintf(w, "No post found for %v", html.EscapeString(req.URL.Path))
    	return
    }

    postData, err := PreparePost(fname)
    if nil != err {
    	logger.Println("Failed getting/parsing content of file:", fname)
    	w.WriteHeader(http.StatusInternalServerError)
    	fmt.Fprintf(w, "Something  went wrong... :(")
    	return
    }
    _, err = io.Copy(w, postData)
    if nil != err {
    	logger.Println("Failed copying content of post: ", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	fmt.Fprintf(w, "Something  went wrong... :(")
    	return
    }
}

// GET "/"
func handleIndexBare(w http.ResponseWriter, req *http.Request) {
    entries, err := os.ReadDir(gconf.Destination)
    if nil != err {
    	logger.Println("Failed reading directory ", entries)
    	w.WriteHeader(http.StatusInternalServerError)
    	fmt.Fprintf(w, "Something  went wrong... :(")
    	return
    }
    for _, entry := range entries {
    	fname := entry.Name()
    	if strings.HasSuffix(fname, ".md") {
    		fmt.Fprintf(w, "<a href='%v'>%v<a></br>", fname, fname)
    	}
    }
}

// GET "/"
func handleIndexFull(w http.ResponseWriter, req *http.Request) {
    WritePreviewsToResponse(w)
}

// GET "/{post}"
func handlePostFull(w http.ResponseWriter, req *http.Request) {

    if req.PathValue("post") == "favicon.ico" {
    	return
    }

    postID := determinePostID(req.PathValue("post"), w)
    if postID < 0 {
    	return
    }

    err := WritePostToResponse(postID, w)
    if nil != err {
    	logger.Println("Failed to write post to response:", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	return
    }
}

// GET "/{post}/{rsrc}"
func handlePostResourceFull(w http.ResponseWriter, req *http.Request) {
    // "resources" is the base resources
    if req.PathValue("post") == "resources" {
    	handleResourcesFull(w, req)
    	return
    }

    postID := determinePostID(req.PathValue("post"), w)
    if postID < 0 {
    	return
    }

    // check if resource type is supported
    fname := req.PathValue("rsrc")
    var ftype string
    if strings.HasSuffix(fname, ".jpeg") {
    	ftype = "jpeg"
    	w.Header().Set("content-type", "image/jpeg")
    } else {
    	logger.Println("Trying to access unknown resource type:", fname)
    	w.WriteHeader(http.StatusNotFound)
    	return
    }

    metadata := postsMetadata[postID]
    fname = filepath.Join(metadata.Path, fname)
    finfo, err := os.Stat(fname)
    if nil != err || !finfo.Mode().IsRegular() {
    	w.WriteHeader(http.StatusNotFound)
    	fmt.Fprintf(w, "No resource found for %v", html.EscapeString(req.URL.Path))
    	return
    }

    resourceData, err := os.ReadFile(fname)
    if nil != err {
    	logger.Println("Failed getting content of resource file:", fname)
    	w.WriteHeader(http.StatusInternalServerError)
    	fmt.Fprintf(w, "Something  went wrong... :(")
    	return
    }

    size, err := io.Copy(w, bytes.NewReader(resourceData))
    if nil != err {
    	logger.Println("Failed copying content of resource: ", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	fmt.Fprintf(w, "Something  went wrong... :(")
    	return
    }

    switch ftype {
    case "jpeg":
    	w.Header().Set("Content-Type", "image/jpeg")
    	w.Header().Set("Content-Length", strconv.Itoa(int(size)))
    default:
    	w.WriteHeader(http.StatusBadRequest)
    	logger.Println("Not Implemented! (2)", ftype)
    }

}

// GET "/resources/{rsrc}"
func handleResourcesFull(w http.ResponseWriter, req *http.Request) {
    fname := req.PathValue("rsrc")
    if strings.HasSuffix(fname, ".css") {
    	w.Header().Set("content-type", "text/css")
    } else if strings.HasSuffix(fname, ".html") {
    	w.Header().Set("content-type", "text/plain")
    } else {
    	w.WriteHeader(http.StatusNotFound)
    	return
    }
    fname = filepath.Join(gconf.Destination, "resources", fname)
    finfo, err := os.Stat(fname)
    if nil != err || !finfo.Mode().IsRegular() {
    	w.WriteHeader(http.StatusNotFound)
    	fmt.Fprintf(w, "No resource found for %v", html.EscapeString(req.URL.Path))
    	return
    }

    resourceData, err := os.ReadFile(fname)
    if nil != err {
    	logger.Println("Failed getting content of resource file:", fname)
    	w.WriteHeader(http.StatusInternalServerError)
    	fmt.Fprintf(w, "Something  went wrong... :(")
    	return
    }
    _, err = io.Copy(w, bytes.NewReader(resourceData))
    if nil != err {
    	logger.Println("Failed copying content of resource: ", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	fmt.Fprintf(w, "Something  went wrong... :(")
    	return
    }
}

func determinePostID(reqId string, w http.ResponseWriter) int {
    id, err := strconv.Atoi(reqId)
    if err != nil {
    	logger.Println("Failed to parse post ID:", err)
    	w.WriteHeader(http.StatusBadRequest)
    	return -1
    }
    if id < 0 || id >= len(postsMetadata) {
    	w.WriteHeader(http.StatusNotFound)
    	return -1
    }
    return id
}
