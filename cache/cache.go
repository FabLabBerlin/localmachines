package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
)

func isSafe(path string) bool {
	return true
}

func forwardHeader(w http.ResponseWriter, resp *http.Response, key string) {
	value := resp.Header[key]
	if len(value) == 1 {
		w.Header().Set(key, value[0])
	}
}

func main() {
	masterHost := flag.String("master", "127.0.0.1:8080", "Master host")
	listenHost := flag.String("listen", "127.0.0.1:1234", "Listen host")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "machines/", 302)
		} else if r.URL.Path == "/api/users/login" {
			form := url.Values{
				"username": {r.FormValue("username")},
				"password": {r.FormValue("password")},
			}
			resp, err := http.PostForm("http://"+*masterHost+"/api/users/login", form)
			if err != nil {
				log.Fatalf("http.Get: %v", err)
			}
			defer resp.Body.Close()
			forwardHeader(w, resp, "Content-Type")
			forwardHeader(w, resp, "Cookie")
			forwardHeader(w, resp, "Set-Cookie")
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				log.Fatalf("io.Copy: %v", err)
			}
		} else if isSafe(r.URL.Path) {
			client := &http.Client{}
			req, err := http.NewRequest("GET", "http://"+*masterHost+r.URL.Path, nil)
			if err != nil {
				log.Fatalf("http.NewRequest: %v", err)
			}
			for _, cookie := range r.Cookies() {
				req.AddCookie(cookie)
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
				log.Fatalf("io.Copy: %v", err)
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
