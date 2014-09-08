package main

import (
	"os"
	"strings"
	"time"

	"github.com/wricardo/batcher"
	"github.com/wricardo/batcher/flushers/function_flusher"
)

var buffer_size int = 10
var buffer_timeout time.Duration = time.Duration(time.Second)

func main() {
	l := NewLogger("output.txt")
	defer l.collector.Shutdown()
	l.Log("aaa")
	l.Log("bbb")
}

type Logger struct {
	filename  string
	collector batcher.Collector
}

func NewLogger(filename string) *Logger {
	l := new(Logger)
	l.filename = filename
	l.collector = createColletor(l.Flush)
	return l
}

func (this *Logger) Log(text string) {
	this.collector.Collect(batcher.String(text))
}

func (this *Logger) Flush(to_flush batcher.Flushable) error {
	if _, err := os.Stat(this.filename); os.IsNotExist(err) {
		os.Create(this.filename)
	}
	f, err := os.OpenFile(this.filename, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString(strings.Join(to_flush.Strings(), "\n") + "\n"); err != nil {
		panic(err)
	}
	return nil
}

func createColletor(flush_func function_flusher.FlushFunction) batcher.Collector {
	f := function_flusher.NewFunctionFlusher(buffer_size, flush_func)
	bcc, _ := batcher.NewBufferedCollectorConfiguration(buffer_size, buffer_timeout)
	return batcher.NewBufferedCollector(bcc, f)
}
