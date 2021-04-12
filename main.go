package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	var trust bool
	var url *url.URL
	port := "8080"
	if hostport, exists := os.LookupEnv("TARGET"); exists {
		if target, err := url.Parse("https://" + hostport); err != nil {
			log.Fatalf("https://%v is not a valid url", hostport)
		} else {
			url = target
		}
	} else {
		log.Fatal("mandatory TARGET env variable not provided")
	}
	if insecure, exists := os.LookupEnv("INSECURE"); exists {
		if insec, err := strconv.ParseBool(insecure); err == nil {
			trust = insec
		}
	}
	if envport, exists := os.LookupEnv("PORT"); exists {
		port = envport
	}
	mim := newManInMiddle(url, trust)
	// insert the "virual" edl API endpoint
	http.HandleFunc("/edl/", mim.edlHandler)
	// hijack the "api" API endpoint
	http.HandleFunc("/api/", mim.apiHandler)
	// reverse proxy everything else
	http.HandleFunc("/", mim.defaultHandler)
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
