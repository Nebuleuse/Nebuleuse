package core

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)
var logFile *os.File

type DashboardLogWriter struct {
}

func (w *DashboardLogWriter) Write(p []byte) (n int, err error) {
	Dispatch("log", p[:len(p)])
	return len(p), nil
}

func initLogging() {
	dash := new(DashboardLogWriter)
	var err error

	logFile, err = os.OpenFile("nebuleuse.log", os.O_RDWR|os.O_APPEND, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	Trace = log.New(ioutil.Discard,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(io.MultiWriter(os.Stdout, dash, logFile),
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(io.MultiWriter(os.Stdout, dash, logFile),
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(os.Stderr, dash, logFile),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func GetPastLogs(size int64) string {
	fi, err := logFile.Stat()
	position := fi.Size() - size
	if position < 0 {
		position = 0
		size = fi.Size()
	}

	_, err = logFile.Seek(position, 0)
	if err != nil {
		Info.Println("Could not seek file", err)
		return ""
	}

	var buffer []byte
	buffer = make([]byte, size)
	_, err = logFile.Read(buffer)

	if err != nil {
		Info.Println(err)
		return ""
	}
	return string(buffer)
}
