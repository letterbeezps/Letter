package testGet

import (
	"letter"
	"testing"
)

type BasicGetResponse struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept          string `json:"Accept"`
		Accept_Encoding string `json:"Accept-Encoding"`
		Accept_Language string `json:"Accept-Language"`
		Host            string `json:"Host"`
		Dnt             string `json:"dst"`
		User_Agent      string `json:"User-Agent"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseNewHeader struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept           string `json:"Accept"`
		Accept_Encoding  string `json:"Accept-Encoding"`
		Accept_Language  string `json:"Accept-Language"`
		Host             string `json:"Host"`
		Dnt              string `json:"dst"`
		User_Agent       string `json:"User-Agent"`
		XWonderfulHeader string `json:"X-Wonder-Header"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseBasicAuth struct {
	Args    struct{} `json:"args"`
	Headers struct {
		Accept          string `json:"Accept"`
		Accept_Encoding string `json:"Accept-Encoding"`
		Accept_Language string `json:"Accept-Language"`
		Host            string `json:"Host"`
		Dnt             string `json:"dst"`
		User_Agent      string `json:"User-Agent"`
		Authorization   string `json:"Authorization"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type BasicGetResponseArgs struct {
	Args struct {
		Goodbye string `json:"Goodbye"`
		Hello   string `json:"Hello"`
	} `json:"args"`
	Headers struct {
		Accept          string `json:"Accept"`
		Accept_Encoding string `json:"Accept-Encoding"`
		Accept_Language string `json:"Accept-Language"`
		Host            string `json:"Host"`
		Dnt             string `json:"dst"`
		User_Agent      string `json:"User-Agent"`
		Authorization   string `json:"Authorization"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

type TestJSONCookies struct {
	Cookies struct {
		AnotherCookie string `json:"AnotherCookie"`
		TestCookie    string `json:"TestCookie"`
	} `json:"cookies"`
}

func TestGetNoOptions(t *testing.T) {
	myJsonStruct := verifyOkResponse(<-letter.GetAsync("http://127.0.0.1:6001/get", nil), t)

	if myJsonStruct.URL != "http://127.0.0.1:6001/get" {
		t.Error("For some reason the URL isn't the same", myJsonStruct.URL)
	}
}

func TestGetSyncNoOptions(t *testing.T) {
	myJsonStruct := verifyOkResponse(letter.Get("http://127.0.0.1:6001/get", nil), t)

	if myJsonStruct.URL != "http://127.0.0.1:6001/get" {
		t.Error("For some reason the URL isn't the same", myJsonStruct.URL)
	}
}

func TestGetNoOptionsChannel(t *testing.T) {
	respChan := letter.GetAsync("http://127.0.0.1:6001/get", nil)
	select {
	case resp := <-respChan:
		verifyOkResponse(resp, t)
	}
}

func TestGetQueryParams(t *testing.T) {
	ro := &letter.RequestOptions{
		Params: map[string]string{"name": "zp"},
	}
	myJsonStruct := verifyOkResponse(letter.Get("http://127.0.0.1:6001/get", ro), t)

	if myJsonStruct.URL != "http://127.0.0.1:6001/get?name=zp" {
		t.Error("For some reason the URL isn't the same", myJsonStruct.URL)
	}
}

func TestGetNoOptionsGzip(t *testing.T) {
	verifyOkResponse(<-letter.GetAsync("http://127.0.0.1:6001/gzip", nil), t)
}

func TestGetWithCookies(t *testing.T) {
	resp := <-letter.GetAsync("http://127.0.0.1:6001/cookies", &letter.RequestOptions{
		Cookies: map[string]string{"AnotherCookie": "Random value", "TestCookie": "Some value"},
	})
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}
	if resp.Ok != true {
		t.Error("Request did not return Ok")
	}
	myJsonStruct := &TestJSONCookies{}
	if err := resp.JSON(myJsonStruct); err != nil {
		t.Error("can not serialize cookie JSON: ", err)
	}
	if myJsonStruct.Cookies.TestCookie != "Some value" {
		t.Errorf("cookie value not is set properly: %v", myJsonStruct)
	}
	if myJsonStruct.Cookies.AnotherCookie != "Random value" {
		t.Errorf("cookie value not is set properly: %v", myJsonStruct)
	}
}

func verifyOkResponse(resp *letter.Response, t *testing.T) *BasicGetResponse {
	if resp.Error != nil {
		t.Error("Unable to make request", resp.Error)
	}

	if resp.Ok != true {
		t.Error("Request did not return OK")
	}

	myJsonStruct := &BasicGetResponse{}
	err := resp.JSON(myJsonStruct)
	if err != nil {
		t.Error("Unable to parse JSON: ", err)
	}

	if myJsonStruct.Headers.Host != "127.0.0.1:6001" {
		t.Error("The host header is invalid")
	}

	return myJsonStruct
}
