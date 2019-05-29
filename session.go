package direwolf

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Session is the main object in direwolf. This is its main features:
// 1. handling redirects
// 2. automatically managing cookies
type Session struct {
	client *http.Client
}

// prepareRequest is to process the parameters from user input.Generate PreRequest object.
func (session Session) prepareRequest(method string, URL string, args ...interface{}) *Request {
	req := new(Request)
	req.Method = strings.ToUpper(method) // Upper the method string
	req.URL = URL

	// Check the type of the paramter and handle it.
	for _, arg := range args {
		switch a := arg.(type) {
		case Headers:
			req.setHeader(a)
		case http.Header:
			req.Headers = a
		case Params:
			req.setParams(a)
		case DataForm:
			req.DataForm = url.Values(a)
		case Data:
			req.Data = a
		case Cookies:
			req.setCookies(a)
		case Proxy:
			req.Proxy = string(a)
		}
	}
	return req
}

// Request is a generic request method.
func (session *Session) request(method string, URL string, args ...interface{}) *Response {
	preq := session.prepareRequest(method, URL, args...)
	return session.send(preq)
}

// Get is a get method.
func (session *Session) Get(URL string, args ...interface{}) *Response {
	return session.request("GET", URL, args...)
}

// Post is a post method.
func (session *Session) Post(URL string, args ...interface{}) *Response {
	return session.request("POST", URL, args...)
}

// send is responsible for handling some subsequent processing of the PreRequest.
func (session *Session) send(preq *Request) *Response {
	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	session.client = &http.Client{
		Transport: trans,
	}
	// New Request
	req, err := http.NewRequest(preq.Method, preq.URL, nil)
	if err != nil {
		panic(err)
	}

	// Add proxy method to transport
	if preq.Proxy != "" {
		proxyURL, err := url.Parse(preq.Proxy)
		if err != nil {
			panic("proxy url has problem")
		}
		trans.Proxy = http.ProxyURL(proxyURL)
	}

	// Handle the Headers.
	req.Header = preq.Headers
	// Handle the DataForm, convert DataForm to strings.Reader.
	// add two new headers: Content-Type and ContentLength.
	if preq.DataForm != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		data := preq.DataForm.Encode()
		req.Body = ioutil.NopCloser(strings.NewReader(data))
		req.ContentLength = int64(len(data))
	}
	// Handle Cookies
	if preq.Cookies != nil {
		for _, cookie := range preq.Cookies {
			req.AddCookie(cookie)
		}
	}

	resp, err := session.client.Do(req) // do request
	if err != nil {
		panic(err)
	}

	buildedResponse := session.buildResponse(preq, resp)

	// build response
	return buildedResponse
}

// buildResponse build response with http.Response after do request.
func (session *Session) buildResponse(req *Request, resp *http.Response) *Response {
	return &Response{
		URL:        req.URL,
		StatusCode: resp.StatusCode,
		Proto:      resp.Proto,
		body:       resp.Body,
	}
}