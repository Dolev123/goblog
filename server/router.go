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
    mux.HandleFunc("GET /", handleIndex)
    mux.HandleFunc("GET /{post}", handlePost)
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
func handlePost(w http.ResponseWriter, req *http.Request) {
    fname := req.PathValue("post")
    if !strings.HasSuffix(fname, ".md") {
	http.Redirect(w, req, "/", http.StatusFound)
	return
    }
    fname = gconf.Destination + fname
    freader, err := os.Open(fname)
    if nil != err {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "No post found for %v", html.EscapeString(req.URL.Path))
	return
    }

    _, err = io.Copy(w, freader)
    if nil != err {
	// TODO:: check if it realy works...
	logger.Println("Failed getting content of file `%v`", fname)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Something  went wrong... :(")
	return
    }
}

// GET "/"
func handleIndex(w http.ResponseWriter, req *http.Request) {
    entries, err := os.ReadDir(gconf.Destination)
    if nil != err {
	logger.Println("Failed reading directory `%v`", entries)
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
