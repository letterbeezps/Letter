package testpost

import (
	"letter"
	"testing"
)

type BasicPostResponse struct {
	Args  struct{} `json:"args"`
	Data  string   `json:"data"`
	Files struct{} `json:"files"`
	Form  struct {
		Name string `json:"name"`
	} `json:"form"`
	Headers struct {
		Accept          string `json:"Accept"`
		Accept_Encoding string `json:"Accept-Encoding"`
		Accept_Language string `json:"Accept-Language"`
		ContentLength   string `json:"Content-Length"`
		ContentType     string `json:"Content-Type"`
		Host            string `json:"Host"`
		User_Agent      string `json:"User-Agent"`
	} `json:"headers"`
	JSON   interface{} `json:"json"`
	Origin string      `json:"origin"`
	URL    string      `json:"url"`
}

type BasicJSONPostResponse struct {
	Args  struct{} `json:"args"`
	Data  string   `json:"data"`
	Files struct{} `json:"files"`
	Form  struct {
		Name string `json:"name"`
	} `json:"form"`
	Headers struct {
		Accept          string `json:"Accept"`
		Accept_Encoding string `json:"Accept-Encoding"`
		Accept_Language string `json:"Accept-Language"`
		ContentLength   string `json:"Content-Length"`
		ContentType     string `json:"Content-Type"`
		Host            string `json:"Host"`
		User_Agent      string `json:"User-Agent"`
		XRequestedWith  string `json:"X-Requested-With"`
	} `json:"headers"`
	JSON struct {
		Name         string        `json:"name"`
		RelationShip letter.Person `json:"relationship"`
	} `json:"json"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicPostFileUpload struct {
	Args  struct{} `json:"args"`
	Data  string   `json:"data"`
	Files struct {
		File string `json:"file"`
	} `json:"files"`
	Form struct {
		Name string `json:"name"`
	} `json:"form"`
	Headers struct {
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		Accept_anguage string `json:"Accept-Language"`
		ContentLength  string `json:"Content-Length"`
		ContentType    string `json:"Content-Type"`
		Host           string `json:"Host"`
		User_Agent     string `json:"User-Agent"`
		XRequestedWith string `json:"X-Requested-With"`
	} `json:"headers"`
	JSON   struct{} `json:"json"`
	Origin string   `json:"origin"`
	URL    string   `json:"url"`
}

type XMLPostMessage struct {
	Name   string
	Age    int
	Height int
}

func TestBasicDataPostRequest(t *testing.T) {
	resp := letter.Post("http://127.0.0.1:6001/post",
		&letter.RequestOptions{
			Data: map[string]string{"name": "zp"},
		})
	verifyOkPostResponse(resp, t)
}

func TestBasicDataPostAsuncRequest(t *testing.T) {
	resp := <-letter.PostAsync("http://127.0.0.1:6001/post",
		&letter.RequestOptions{
			Data: map[string]string{"name": "zp"},
		})
	verifyOkPostResponse(resp, t)
}

func TestBasicDataPostInvalidURLRequest(t *testing.T) {
	resp := <-letter.PostAsync("%../dir/",
		&letter.RequestOptions{
			Data: map[string]string{"name": "zp"},
		})
	verifyOkPostResponse(resp, t)
}

func TestBasicXMLRequest(t *testing.T) {
	xmlDate := XMLPostMessage{
		Name:   "zp",
		Age:    24,
		Height: 176,
	}
	resp := letter.Post("http://127.0.0.1:6001/post",
		&letter.RequestOptions{XML: xmlDate, IsAjax: true})
	myJsonStruct := &BasicJSONPostResponse{}
	resp.JSON(myJsonStruct)
}

func TestBasicJSONRequest(t *testing.T) {
	pData := letter.PostData{
		Name: "zp",
		RelationShip: letter.Person{
			Name: "dd",
		},
	}
	resp := letter.Post("http://127.0.0.1:6001/post",
		&letter.RequestOptions{JSON: pData, IsAjax: true})
	myJsonStruct := &BasicJSONPostResponse{}
	resp.JSON(myJsonStruct)
}

func TestBasicPostRequestUpload(t *testing.T) {
	fd, err := letter.FileUploadFromDisk("upload.txt")
	if err != nil {
		t.Error("Unable to open file: ", err)
	}

	resp := letter.Post("http://127.0.0.1:6001/post",
		&letter.RequestOptions{
			File: fd,
			Data: map[string]string{"name": "zp"},
		})
	if resp.Error != nil {
		t.Fatal("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("request did not return OK")
	}

	myJsonStruct := &BasicPostFileUpload{}
	resp.JSON(myJsonStruct)
}

func verifyOkPostResponse(resp *letter.Response, t *testing.T) *BasicPostResponse {
	if resp.Error != nil {
		t.Fatal("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("request did not return OK")
	}

	myJsonStruct := &BasicPostResponse{}
	// resp.JSON(myJsonStruct)

	if err := resp.JSON(myJsonStruct); err != nil {
		t.Error("Unable to coerce to JSON", err)
	}

	if myJsonStruct.URL != "http://127.0.0.1:6001/post" {
		t.Errorf("the URL isn't the same: %v", myJsonStruct.URL)
	}

	if myJsonStruct.Headers.Host != "127.0.0.1:6001" {
		t.Error("The host header is invalid")
	}

	if myJsonStruct.Form.Name != "zp" {
		t.Errorf("Invalid post response: %#v", myJsonStruct.Form)
	}

	return myJsonStruct
}
