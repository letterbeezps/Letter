package example

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestFast(t *testing.T) {
	data := Data{
		Name: "zp",
		Age:  12,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	client := &fasthttp.Client{}
	// req := &fasthttp.Request{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:6001/post")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	req.SetBody(bytes)
	req.PostArgs().Add("name", "zp")

	resp := &fasthttp.Response{}

	err = client.Do(req, resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp.Body()))
	fmt.Println(resp.ConnectionClose())
	fmt.Println(resp.StatusCode())
}
