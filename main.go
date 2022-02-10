package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	client := http.Client{}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Request received (%v), sending it to the origin ..", req)

		// Handle health checker
		if strings.HasPrefix(req.UserAgent(), "ELB-HealthChecker") {
			w.WriteHeader(200)
			w.Write([]byte("OK\n"))
			return
		}

		// Reverse proxy the request to the origin
		switch req.Host {
		case "secure-api.iamport.kr":
			req.Host = "api.iamport.kr"
		case "secure-service.iamport.kr":
			req.Host = "service.iamport.kr"
		default:
			log.Printf("Denied unallowed Host: %v", req.Host)
			w.WriteHeader(403)
			return
		}
		req.RequestURI = ""
		req.URL.Scheme = "https"
		req.URL.Host = req.Host
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Error while sending request to the origin: %v", err)
			w.WriteHeader(500)
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
			if _, exists := predefined[key]; exists {
				continue
			}
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		// Copy the status
		w.WriteHeader(res.StatusCode)
		// Copy the body
		_, err = io.Copy(w, res.Body)
		if err != nil {
			log.Printf("Error while copying response body: %v", err)
		}

		log.Printf("Successfully proxied response (%v)", res)
	})

	log.Printf("Starting server on :80 ...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
