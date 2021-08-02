package letter

import (
	"github.com/valyala/fasthttp"
)

func Get(url string, ro *RequestOptions) *Response {
	return doRequest(fasthttp.MethodGet, url, ro)
}

func GetAsync(url string, ro *RequestOptions) chan *Response {
	return doAsyncRequest(fasthttp.MethodGet, url, ro)
}

func Post(url string, ro *RequestOptions) *Response {
	return doRequest(fasthttp.MethodPost, url, ro)
}

func PostAsync(url string, ro *RequestOptions) chan *Response {
	return doAsyncRequest(fasthttp.MethodPost, url, ro)
}
