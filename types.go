package httputils

import "net/http"

type RespChannel chan *Response

type Response struct {
	rChannel chan *http.Response
	r *http.Response
	current string
	body []byte
	err error
}

type HTTPError struct{
	Code int
	Label string
}