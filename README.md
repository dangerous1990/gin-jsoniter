<img src="https://travis-ci.org/dangerous1990/gin-jsoniter.svg?branch=master">

# Installation
```
go get -u github.com/dangerous1990/gin-jsoniter
```
go 1.11.1+  add in your go.mod file
```
go mod edit --require github.com/dangerous1990/gin-jsoniter@v1.0.0
```
# Quick start
```go
    type Form struct {
    	Name string `json:"name"`
    }
    
    c.ShouldBindWith(&form, jsoniterserialize.JsoniterBinding{jsoniterAPI})

    c.Render(http.StatusOK, jsoniterserialize.jsoniterRender{jsoniterAPI, map[string]string{"message": "hello " + form.Name}})
```

# More information please see gin and jsoniter doc
* [gin](https://github.com/gin-gonic/gin)
* [jsoniter](https://github.com/json-iterator)