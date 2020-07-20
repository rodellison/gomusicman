package mocks

import (
	"net/http"
)

type MockHTTPClient struct {
	//This DoFunc func (with the same signature as the httpclient's Do function allows us to ensure the struct
	//matches the signature of httpclient, while also separately allowing us to define our own Do function.
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	//GetDoFunc is a value that gets set during test, that allows providing specific canned responses or errors
	GetDoHTTPFunc func(req *http.Request) (*http.Response, error)
	/*e.g.
	return &http.Response{
	  StatusCode: 200,
	  Body:       <whatever string or json we want>,
	}, nil
	*/
)

// Do is the mock client's `Do` func
// By establishing this 'fake', we can have it call a function represented by the GetDoFunc variable.
//This main thing this does is allow us to define an anonymous function that returns whatever canned response
//we want for a given test's call
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoHTTPFunc(req)
}
