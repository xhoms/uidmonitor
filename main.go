package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var hostport, port, insecure string
	var exists, trust bool
	if hostport, exists = os.LookupEnv("TARGET"); !exists {
		log.Println("Mandatory env variable TARGET missing")
		os.Exit(1)
	}
	if insecure, exists = os.LookupEnv("INSECURE"); exists {
		if insec, err := strconv.ParseBool(insecure); err == nil {
			trust = insec
		}
	}
	if port, exists = os.LookupEnv("PORT"); !exists {
		port = "8080"
	}
	if mim, err := newMim(hostport, trust); err == nil {
		// insert the "virual" edl API endpoint
		http.HandleFunc("/edl/", mim.edl)
		// hijack the "api" API endpoint
		http.HandleFunc("/api/", mim.proxyUID)
		// reverse proxy everything else
		http.HandleFunc("/", mim.proxyReverse)
		log.Println("Listening on port", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(err)
	}
}
