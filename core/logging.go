package core

import (
	"io"
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

	normalOut := io.MultiWriter(os.Stdout, dash, logFile)
	errOut := io.MultiWriter(os.Stderr, dash, logFile)

	Trace = log.New(normalOut,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(normalOut,
		"INFO: ",
		log.Ldate|log.Ltime)

	Warning = log.New(normalOut,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errOut,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	log.SetOutput(errOut)
	log.SetPrefix("ERROR: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
		Error.Println("Could not seek file", err)
		return ""
	}

	var buffer []byte
	buffer = make([]byte, size)
	_, err = logFile.Read(buffer)

	if err != nil {
		Error.Println(err)
		return ""
	}
	return string(buffer)
}
