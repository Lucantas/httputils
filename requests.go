package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Fetch(url string, params ...map[string]interface{}) *Response {
	var r Response
	out := make(chan *http.Response)
	go func() {
		c := &http.Client{}

		resp, err := c.Do(newRequest(url, params...))

		if err != nil {
			log.Println(err)
		}
		out <- resp

		close(out)
	}()
	r.rChannel = out
	return &r
}

func (resp *Response) Then() *Response {
	resp.r = <-resp.rChannel
	return resp
}

func (resp *Response) BodyReader() io.ReadCloser {
	return resp.r.Body
}

func (resp *Response) Body() *Response {
	r := resp.r

	if r == nil {
		return resp
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
	}

	resp.body = b
	return resp
}

func (resp *Response) Bytes() []byte {
	return resp.body
}

func (resp *Response) String() string {
	return string(resp.body)
}

func (resp *Response) JSON(v interface{}) (*Response, error) {
	return resp, json.Unmarshal(resp.body, &v)
}

func (resp *Response) Catch() *Response {
	r := resp.r
	if r == nil {
		resp.err = errors.New("URL Inválida")
		return resp
	}
	if r.StatusCode < http.StatusOK || r.StatusCode > 299 {
		resp.err = errors.New(fmt.Sprintf("HTTP Error | Message: %s | Code: %d", r.Status, r.StatusCode))
		return resp
	}
	return nil
}

func applyParams(req *http.Request, p map[string]string) *http.Request {
	for key, value := range p {
		req.Header.Set(key, value)
	}
	return req
}

func newRequest(url string, params ...map[string]interface{}) *http.Request {
	if len(params) > 0 {
		// only the first map string interface will be considered
		return requestWithParams(url, params[0])
	}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println(err)
	}
	return req
}

func requestWithParams(url string, params map[string]interface{}) *http.Request {
	// set method
	var req *http.Request
	var err error
	method := getMethod(params)
	body := getBody(params)
	if body != nil {
		req, err = http.NewRequest(method, url, body)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		log.Println(err)
	}

	req = setHeaders(req, params)
	return req
}

func getBody(p map[string]interface{}) *strings.Reader {
	for key, i := range p {
		switch key {
		case "body":
			if s, ok := i.(string); ok {
				return strings.NewReader(s)
			}

			if bs, ok := i.([]byte); ok {
				return strings.NewReader(string(bs))
			}
		}
	}
	return nil
}

func setHeaders(req *http.Request, p map[string]interface{}) *http.Request {
	for key, i := range p {
		switch key {
		case "headers":
			if headers, ok := i.(map[string]string); ok {
				for k, y := range headers {
					req.Header.Set(k, y)
				}
			}
		}
	}
	return req
}

func getMethod(p map[string]interface{}) string {
	for key, i := range p {
		switch key {
		case "method":
			if s, ok := i.(string); ok {
				return s
			}
		}
	}
	return ""
}
