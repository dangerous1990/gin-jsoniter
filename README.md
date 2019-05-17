* custom gin binding and render
    * custom binding implements binding interface
    * custom render implements render interface
* test ShouldBindWith(obj interface{}, b binding.Binding) and  Render(code int, r render.Render)

# custom gin binding and render

## custom binding implements binding interface
```go
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

```
## custom render implements render interface
```go
var (
	jsonContentType               = []string{"application/json; charset=utf-8"}
	_               render.Render = jsoniterRender{}
)

type jsoniterRender struct {
	jsoniterAPI jsoniter.API
	Data        interface{}
}

func (render jsoniterRender) Render(w http.ResponseWriter) error {
	bytes, err := render.jsoniterAPI.Marshal(render.Data)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

func (render jsoniterRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

```

# test ShouldBindWith(obj interface{}, b binding.Binding) and  Render(code int, r render.Render)
 ```go
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

```