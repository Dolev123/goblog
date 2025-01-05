package server

import (
    "html"
    "fmt"
    "io"
    "net/http"
    "strings"
    "os"

    "github.com/Dolev123/goblog/config"
    pkglog "github.com/Dolev123/goblog/logger"
)


var gconf *config.Config

var logger = pkglog.CreateNewLogger()

func StartServer(conf *config.Config) {
    gconf = conf
    mux := http.NewServeMux()
    switch gconf.Structure {
    case "bare": 
        mux.HandleFunc("GET /", handleIndexBare)
	mux.HandleFunc("GET /{post}", handlePostBare)
    case "full":
	logger.Fatal("full mode not yet implemented!")
        mux.HandleFunc("GET /", handleIndexFull)
	mux.HandleFunc("GET /{post}", handlePostFull)
    default:
	logger.Fatal("Unknown server mode/structure:", gconf.Structure)
    }
    logger.Println("setting server in", gconf.Structure, "mode")
    srv := &http.Server{
	Addr: conf.ListenAddr,
	Handler: mux,
	// should make a new logger for error?
	ErrorLog: logger,
	// TODO:: add TLS
    }
    logger.Fatal("Server Failed with:", srv.ListenAndServe())
}

// GET "/{post}"
func handlePostBare(w http.ResponseWriter, req *http.Request) {
    fname := req.PathValue("post")
    if !strings.HasSuffix(fname, ".md") {
	http.Redirect(w, req, "/", http.StatusFound)
	return
    }
    fname = gconf.Destination + fname
    finfo, err := os.Stat(fname)
    if nil != err || !finfo.Mode().IsRegular() {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "No post found for %v", html.EscapeString(req.URL.Path))
	return
    }

    postData, err := PreparePost(fname)
    if nil != err {
	// TODO:: check if it realy works...
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
    logger.Fatal("Not Implemented!")
}
// GET "/{post}"
func handlePostFull(w http.ResponseWriter, req *http.Request) {
    logger.Fatal("Not Implemented!")
}
