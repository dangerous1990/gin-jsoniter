package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

type Form struct {
	Name string `json:"name"`
}

var jsoniterAPI = jsoniter.Config{}.Froze()

func jsoniterTest(c *gin.Context) {
	var form Form
	if err := c.ShouldBindWith(&form, JsoniterBinding{jsoniterAPI}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	}
	c.Render(http.StatusOK, jsoniterRender{jsoniterAPI, map[string]string{"message": "hello " + form.Name}})
}

func TestCheckMails(t *testing.T) {
	router := gin.Default()
	router.POST("/test", jsoniterTest)
	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"name":"world"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, jsoniterAPI.Get(w.Body.Bytes(), "message").ToString(), "hello world")
}
