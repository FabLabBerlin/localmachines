package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
)

func isSafe(path string) bool {
	return true
}

func main() {
	masterHost := flag.String("master", "127.0.0.1:8080", "Master host")
	listenHost := flag.String("listen", "127.0.0.1:1234", "Listen host")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "machines/", 302)
		} else if isSafe(r.URL.Path) {
			resp, err := http.Get("http://" + *masterHost + r.URL.Path)
			if err != nil {
				log.Fatalf("http.Get: %v", err)
			}
			contentType := resp.Header["Content-Type"]
			if len(contentType) == 1 {
				w.Header().Set("Content-Type", contentType[0])
			}
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
