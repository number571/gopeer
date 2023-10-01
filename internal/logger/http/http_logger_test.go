package http

import (
	"net/http"
	"testing"
)

const (
	tcService = "TST"
	tcFmtLog  = "service=TST method=GET path=/api/index conn= message=hello_world"
)

func TestLogger(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/index", nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.RemoteAddr = "127.0.0.1:55555"

	logBuilder := NewLogBuilder(tcService, req)
	logGetter := logBuilder.WithMessage("hello_world").Get()

	if logGetter.GetConn() != "127.0.0.1:55555" {
		t.Error("got conn != conn")
		return
	}

	if logGetter.GetMessage() != "hello_world" {
		t.Error("got message != message")
		return
	}

	if logGetter.GetMethod() != "GET" {
		t.Error("got method != method")
		return
	}

	if logGetter.GetPath() != "/api/index" {
		t.Error("got path != path")
		return
	}

	if logGetter.GetService() != "TST" {
		t.Error("got service != service")
		return
	}
}
