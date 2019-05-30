package direwolf

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

// Session is the main object in direwolf. This is its main features:
// 1. handling redirects
// 2. automatically managing cookies
type Session struct {
	Client    *http.Client
	Transport *http.Transport
	Cookies   *cookiejar.Jar
	Headers   *http.Header
	Proxy     string
	Timeout   int
}

// prepareRequest is to process the parameters from user input.Generate PreRequest object.
func (session Session) prepareRequest(method string, URL string, args ...interface{}) *RequestSetting {
	req := new(RequestSetting)
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
func (session *Session) Request(method string, URL string, args ...interface{}) *Response {
	preq := session.prepareRequest(method, URL, args...)
	return session.send(preq)
}

// Get is a get method.
func (session *Session) Get(URL string, args ...interface{}) *Response {
	return session.Request("GET", URL, args...)
}

// Post is a post method.
func (session *Session) Post(URL string, args ...interface{}) *Response {
	return session.Request("POST", URL, args...)
}

// send is responsible for handling some subsequent processing of the PreRequest.
func (session *Session) send(preq *RequestSetting) *Response {
	buildedResponse := Download(preq, session.Client, session.Transport)

	// build response
	return buildedResponse
}

// NewSession make a Session, and set a default Client and Transport.
func NewSession() *Session {
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

	// set CookieJar
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		panic("proxy url has problem")
	}

	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	return &Session{
		Client:    client,
		Transport: trans,
	}
}
