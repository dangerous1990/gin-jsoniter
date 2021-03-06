package jsoniterserialize

import (
	"net/http"

	"github.com/gin-gonic/gin/render"
	jsoniter "github.com/json-iterator/go"
)

var (
	jsonContentType               = []string{"application/json; charset=utf-8"}
	_               render.Render = JsoniterRender{}
)

type JsoniterRender struct {
	JsoniterAPI jsoniter.API
	Data        interface{}
}

func (render JsoniterRender) Render(w http.ResponseWriter) error {
	bytes, err := render.JsoniterAPI.Marshal(render.Data)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

func (render JsoniterRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
