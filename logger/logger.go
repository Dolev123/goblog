package logger

import (
    "fmt"
    "log"
    "runtime"
    "strings"
)

func determineLoggerPrefix() string {
    // Stack order:
    // 0: logger.determineLoggerPrefix
    // 1: logger.CreateNewLogger
    // 2: <calling_package>
    pc, _, _, ok := runtime.Caller(2)
    if !ok {
    	log.Fatal("Could not get stack[2] for logging")
    }
    fname := runtime.FuncForPC(pc).Name()
    // package is "full/pkg/path.function" -> only interested in package name
    start := strings.LastIndex(fname, "/")
    end := strings.LastIndex(fname, ".")
    if start+1 < end {
    	fname = fname[start+1 : end]
    }
    return strings.ToLower(fname)
}

func CreateNewLogger() *log.Logger {
    prefix := fmt.Sprintf("%-12s", "["+determineLoggerPrefix()+"]")
    writer := log.Writer()
    flags := log.Ldate | log.Ltime | log.Lshortfile
    return log.New(writer, prefix, flags)
}
