package letter

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type Response struct {
	Ok          bool
	Error       error
	StatusCode  int
	Header      fasthttp.ResponseHeader
	RawResponse *fasthttp.Response
}

func buildResponse(resp *fasthttp.Response, err error) *Response {
	if err != nil {
		return &Response{
			Error: err,
		}
	}

	resp.SetConnectionClose()

	return &Response{
		Ok:          resp.StatusCode() >= 200 && resp.StatusCode() < 300,
		Error:       nil,
		StatusCode:  resp.StatusCode(),
		Header:      resp.Header,
		RawResponse: resp,
	}
}

func (r *Response) Close() error {
	if r.Error != nil {
		return fmt.Errorf("Unable to make request %s", r.Error)
	}
	r.RawResponse.SetConnectionClose()
	return nil
}

func (r *Response) JSON(userStruct interface{}) error {
	// fmt.Println(r.RawResponse.Body())
	fmt.Println(string(r.RawResponse.Body()))
	if err := json.Unmarshal(r.RawResponse.Body(), &userStruct); err != nil {
		return err
	}
	return nil
}
