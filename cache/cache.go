package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
)

var masterHost *string

func isSafe(path string) bool {
	return true
}

func forwardHeader(w http.ResponseWriter, resp *http.Response, key string) {
	value := resp.Header[key]
	if len(value) == 1 {
		w.Header().Set(key, value[0])
	}
}

func newRequest(browserReq *http.Request) (req *http.Request, err error) {
	var body io.Reader
	if browserReq.Method != "GET" {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, browserReq.Body); err != nil {
			return nil, fmt.Errorf("io.Copy r.Body: %v", err)
		}
		body = buf
	}
	uri := "http://" + *masterHost + browserReq.URL.Path
	req, err = http.NewRequest(browserReq.Method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}
	for _, cookie := range browserReq.Cookies() {
		req.AddCookie(cookie)
	}
	if len(browserReq.Header["Content-Type"]) == 1 {
		req.Header.Set("Content-Type", browserReq.Header["Content-Type"][0])
	}
	return
}

func main() {
	masterHost = flag.String("master", "127.0.0.1:8080", "Master host")
	listenHost := flag.String("listen", "127.0.0.1:1234", "Listen host")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "machines/", 302)
		} else if isSafe(r.URL.Path) {
			client := &http.Client{}
			req, err := newRequest(r)
			if err != nil {
				log.Fatalf("http.NewRequest: %v", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("http.Get: %v", err)
			}
			defer resp.Body.Close()
			forwardHeader(w, resp, "Content-Type")
			forwardHeader(w, resp, "Cookie")
			forwardHeader(w, resp, "Set-Cookie")
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				log.Fatalf("io.Copy resp.Body: %v", err)
			}
		} else {
			fmt.Fprintf(w, "Sorry not possible. Thank you.")
		}
	})

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(*listenHost, nil))
}
