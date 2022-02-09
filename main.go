package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	client := http.Client{}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Request received (%v), sending it to the origin ..", req)

		// Reverse proxy the request to the origin
		req.RequestURI = ""
		req.URL.Scheme = "https"
		req.URL.Host = req.Host
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error while sending request to the origin: %v", err)
			return
		}

		// Proxy the response back to the client
		predefined := map[string]string{
			"Strict-Transport-Security": "max-age=15552000",
			"X-Content-Type-Options":    "nosniff",
			"X-Frame-Options":           "SAMEORIGIN",
			"X-Xss-Protection":          "1; mode=block",
		}
		for key, value := range predefined {
			w.Header().Add(key, value)
		}
		// Copy the headers
		for key, values := range res.Header {
			if _, exists := predefined["foo"]; exists {
				continue
			}
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		// Copy the status
		w.WriteHeader(res.StatusCode)
		// Copy the body
		io.Copy(w, res.Body)

		log.Printf("Successfully proxied response (%v)", res)
	})

	log.Printf("Starting server on :80 ...")
	http.ListenAndServe(":80", nil)
}
