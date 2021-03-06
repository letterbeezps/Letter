package letter

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"

	"github.com/valyala/fasthttp"
)

type RequestOptions struct {

	// a map of key values that will eventually convert into the query string of
	// a GET string or a POST request
	Data map[string]string

	// GET request
	Params map[string]string

	// JSON
	JSON interface{}

	// XML
	XML interface{}

	// FILE
	File *FileUpload

	// Headers
	Headers map[string]string

	// IsAjax is a flag that can be set to make request appear to be generated by browser javascript
	IsAjax bool

	// Cookie
	// Cookies []fasthttp.Cookie
	Cookies map[string]string
}

func doRequest(httpMethod, url string, ro *RequestOptions) *Response {
	return buildResponse(buildRequest(httpMethod, url, ro))
}

func doAsyncRequest(httpMethod, url string, ro *RequestOptions) chan *Response {
	responseChan := make(chan *Response, 1)
	go func() {
		responseChan <- doRequest(httpMethod, url, ro)
	}()
	return responseChan
}

func buildRequest(httpMethod, userURL string, ro *RequestOptions) (*fasthttp.Response, error) {
	if ro == nil {
		ro = &RequestOptions{}
	}
	httpCLient := buildHTTPClient(ro)

	var err error

	req, err := buildHTTPRequest(httpMethod, userURL, ro)
	if err != nil {
		return nil, err
	}
	// defer fasthttp.ReleaseRequest(req)

	if ro.Params != nil {
		err = buildURLParams(req, ro.Params)
		if err != nil {
			panic(err)
		}
	}

	resp := fasthttp.AcquireResponse()
	// defer fasthttp.ReleaseResponse(resp)

	addHeaders(ro, req)

	addCookies(ro, req)

	httpCLient.Do(req, resp)

	return resp, nil

}

func buildHTTPRequest(httpMethod, userURL string, ro *RequestOptions) (*fasthttp.Request, error) {

	if ro.JSON != nil {
		return createJSONRequest(httpMethod, userURL, ro)
	}

	if ro.XML != nil {
		return createBasicXMLRequest(httpMethod, userURL, ro)
	}

	if ro.File != nil {
		return createMultiPartPostRequest(httpMethod, userURL, ro)
	}

	if ro.Data != nil {
		return createDataRequest(httpMethod, userURL, ro)
	}

	return createBasicRequest(httpMethod, userURL, nil)
}

func createBasicRequest(httpMethod, userURL string, ro *RequestOptions) (*fasthttp.Request, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(userURL)
	req.Header.SetMethod(httpMethod)
	return req, nil
}

func createDataRequest(httpMethod, userURL string, ro *RequestOptions) (*fasthttp.Request, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(userURL)
	req.Header.SetMethod(httpMethod)
	req.Header.SetContentType("application/x-www-form-urlencoded")

	for k, v := range ro.Data {
		req.PostArgs().Add(k, v)
	}

	return req, nil
}

func createJSONRequest(httpMethod, userURL string, ro *RequestOptions) (*fasthttp.Request, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(userURL)
	req.Header.SetMethod(httpMethod)
	req.Header.SetContentType("application/json")

	tempBuffer := &bytes.Buffer{}
	if err := json.NewEncoder(tempBuffer).Encode(ro.JSON); err != nil {
		return nil, err
	}

	req.SetBody(tempBuffer.Bytes())

	return req, nil
}

func createMultiPartPostRequest(httpMethod, userURL string, ro *RequestOptions) (*fasthttp.Request, error) {
	if httpMethod != fasthttp.MethodPost {
		return nil, errors.New("Only POST method is valid for MultiPartPostRequest")
	}
	requestBody := &bytes.Buffer{}

	multipartWriter := multipart.NewWriter(requestBody)
	writer, err := multipartWriter.CreateFormFile("file", ro.File.FileName)
	if err != nil {
		return nil, err
	}

	if ro.File.FileContents == nil {
		return nil, errors.New("pointer FileContents cannot be nil")
	}

	defer ro.File.FileContents.Close()

	if _, err = io.Copy(writer, ro.File.FileContents); err != nil && err != io.EOF {
		return nil, err
	}

	for k, v := range ro.Data {
		multipartWriter.WriteField(k, v)
	}

	// stop write data
	if err = multipartWriter.Close(); err != nil {
		return nil, err
	}
	contentType := multipartWriter.FormDataContentType()
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(userURL)
	req.Header.SetMethod(httpMethod)
	req.Header.SetContentType(contentType)
	req.SetBody(requestBody.Bytes())

	return req, nil
}

func createBasicXMLRequest(httpMethod, userURL string, ro *RequestOptions) (*fasthttp.Request, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(userURL)
	req.Header.SetMethod(httpMethod)
	req.Header.SetContentType("application/xml")

	tempBuffer := &bytes.Buffer{}
	if err := json.NewEncoder(tempBuffer).Encode(ro.XML); err != nil {
		return nil, err
	}

	req.SetBody(tempBuffer.Bytes())

	return req, nil
}

func buildHTTPClient(ro *RequestOptions) *fasthttp.Client {
	return &fasthttp.Client{
		Name: "letterClient",
	}
}

func buildURLParams(req *fasthttp.Request, params map[string]string) error {
	for key, value := range params {
		req.URI().QueryArgs().Add(key, value)
	}
	return nil
}

func addHeaders(ro *RequestOptions, req *fasthttp.Request) {
	for k, v := range ro.Headers {
		req.Header.Add(k, v)
	}
	if ro.IsAjax {
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
	}
}

func addCookies(ro *RequestOptions, req *fasthttp.Request) {
	for k, v := range ro.Cookies {
		req.Header.SetCookie(k, v)
	}
}
