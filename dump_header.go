package main;

import (
    "net/http"
    "fmt"
    "html"
    "log"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request URL: %q\n", html.EscapeString(r.URL.Path))
	for k,v := range r.Header {
	    fmt.Fprintf(w, "%q : %q\n", k, v)
	}
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
