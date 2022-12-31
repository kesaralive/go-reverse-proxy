package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {

	// define origin server URL
	originServerURL, err := url.Parse("http://127.0.0.1:8081")
	if err != nil {
		log.Fatal(err)
	}

	reverseProxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[reverse proxy server] recieved request at: %s\n", time.Now())

		// set req Host, URL and Request URI to forward a request to the origin server
		r.Host = originServerURL.Host
		r.URL.Host = originServerURL.Host
		r.URL.Scheme = originServerURL.Scheme
		r.RequestURI = ""

		// send a request to the origin server
		originServerResponse, err := http.DefaultClient.Do(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, err)
			return
		}

		// return response to the client
		w.WriteHeader(http.StatusOK)
		io.Copy(w, originServerResponse.Body)
	})
	log.Fatal(http.ListenAndServe(":8080", reverseProxy))
}
