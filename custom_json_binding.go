package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	jsoniter "github.com/json-iterator/go"
)

var (
	_ binding.Binding     = JsoniterBinding{}
	_ binding.BindingBody = JsoniterBinding{}
)

type JsoniterBinding struct {
	jsoniterAPI jsoniter.API
}

func (jsoniterBinding JsoniterBinding) Name() string {
	return "json"
}

func (jsoniterBinding JsoniterBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeJSON(jsoniterBinding, req.Body, obj)
}

func (jsoniterBinding JsoniterBinding) BindBody(body []byte, obj interface{}) error {
	return decodeJSON(jsoniterBinding, bytes.NewReader(body), obj)
}

func decodeJSON(jsoniterBinding JsoniterBinding, r io.Reader, obj interface{}) error {
	decoder := jsoniterBinding.jsoniterAPI.NewDecoder(r)
	if binding.EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}

func validate(obj interface{}) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
