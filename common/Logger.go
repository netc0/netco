package common

import (
	"log"
	"io"
	"fmt"
	"bytes"
)

type xLoggerWriter struct {
	io.Writer
}

func (this* xLoggerWriter)Write(p []byte) (n int, err error) {
	size := len(p)
	print(string(p))
	return size, nil
}

type ILogger interface {
	Prefix(v string)
	Debug(v ... interface{})
}

type dLogger struct {
	ILogger
	log.Logger
}

// default logger impl
func (this *dLogger) buildString(obj ... interface{}) []byte {
	var buffer bytes.Buffer
	flag := false
	for _, v := range obj {
		s := fmt.Sprintf("%v ", v)
		buffer.WriteString(s)
		flag = true
	}
	if !flag {
		return nil
	}

	return buffer.Bytes()
}

func (this* dLogger)Prefix(v string) {
	this.SetPrefix(v)
}

func (this* dLogger) Debug(v ... interface{}) {
	s := this.buildString(v...)
	if s == nil { return }
	this.Output(2, string(s))
}

func GetLogger () ILogger{
	var l dLogger
	l.SetFlags(log.LstdFlags | log.Lshortfile)
	l.SetOutput(&xLoggerWriter{})
	return &l
}
