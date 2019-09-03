package direwolf

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMini(t *testing.T) {
	session := NewSession()
	proxy := &Proxy{
		HTTP:  "http://127.0.0.1:12333",
		HTTPS: "http://127.0.0.1:12333",
	}
	session.Proxy = proxy
	resp, err := session.Get("http://httpbin.org/ip")
	if err != nil {
		return
	}
	fmt.Println(resp.Text())
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is GET before going to check other features
		if r.Method != "GET" {
			t.Fatalf("Expected method %q; got %q", "GET", r.Method)
		}
		if r.URL.Path == "/test" {
			w.Write([]byte("passed"))
		}
	}))
	defer ts.Close()

	resp, err := Get(ts.URL + "/test")
	if err != nil {
		t.Fatal(err)
	}
	text := resp.Text()
	if text != "passed" {
		t.Fatal("response was wrong, not", text)
	}
	t.Log("Get test passed")
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is POST before going to check other features
		if r.Method != "POST" {
			t.Fatalf("Expected method %q; got %q", "POST", r.Method)
		}
		if r.URL.Path == "/test" {
			w.Write([]byte("passed"))
		}
		body, _ := ioutil.ReadAll(r.Body)
		bodyString := string(body)
		if bodyString != "direwolf" {
			t.Fatal("Request body was wrong, not", bodyString)
		}
	}))
	defer ts.Close()

	body := Body("direwolf")
	resp, err := Post(ts.URL+"/test", body)
	if err != nil {
		t.Fatal(err)
	}
	text := resp.Text()
	if text != "passed" {
		t.Fatal("response was wrong, not", text)
	}
	t.Log("Post test passed")
}

func TestRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is GET before going to check other features
		if r.Method != "GET" {
			t.Fatalf("Expected method %q; got %q", "GET", r.Method)
		}
		if r.URL.Path == "/test" {
			w.Write([]byte("passed"))
		}
	}))
	defer ts.Close()

	req := NewRequestSetting("Get", ts.URL+"/test")
	resp, err := Request(req)
	if err != nil {
		t.Fatal(err)
	}
	text := resp.Text()
	if text != "passed" {
		t.Fatal("response was wrong, not", text)
	}
	t.Log("Request test passed")
}

func TestSendCookie(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is GET before going to check other features
		if r.Method != "GET" {
			t.Fatalf("Expected method %q; got %q", "GET", r.Method)
		}
		if r.URL.Path == "/test" {
			w.Write([]byte("passed"))
		}
		cookie, err := r.Cookie("name")
		if err != nil {
			t.Fatal(err)
		}
		if cookie.Value != "direwolf" {
			t.Fatalf("Expected value %q; got %q", "direwolf", cookie)
		}
	}))
	defer ts.Close()

	cookies := NewCookies(
		"name", "direwolf",
	)
	resp, err := Get(ts.URL+"/test", cookies)
	if err != nil {
		t.Fatal(err)
	}
	text := resp.Text()
	if text != "passed" {
		t.Fatal("response was wrong, not", text)
	}
	t.Log("request cookies test passed")
}

func TestRequestHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is GET before going to check other features
		if r.Method != "GET" {
			t.Fatalf("Expected method %q; got %q", "GET", r.Method)
		}
		if r.URL.Path == "/test" {
			if r.Header.Get("User-Agent") == "direwolf" {
				w.Write([]byte("passed"))
			}
		}
	}))
	defer ts.Close()

	headers := NewHeaders(
		"User-Agent", "direwolf",
	)
	resp, err := Get(ts.URL+"/test", headers)
	if err != nil {
		t.Fatal(err)
	}
	text := resp.Text()
	if text != "passed" {
		t.Fatal("response was wrong, not", text)
	}
	t.Log("request headers test passed")
}

func TestParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method is GET before going to check other features
		if r.Method != "GET" {
			t.Fatalf("Expected method %q; got %q", "GET", r.Method)
		}
		if r.URL.Path == "/test" {
			if r.FormValue("key") == "value" {
				w.Write([]byte("passed"))
			}
		}
	}))
	defer ts.Close()

	params := NewParams(
		"key", "value",
	)
	resp, err := Get(ts.URL+"/test", params)
	if err != nil {
		t.Fatal(err)
	}
	text := resp.Text()
	if text != "passed" {
		t.Fatal("response was wrong, not", text)
	}
	t.Log("request headers test passed")
}
