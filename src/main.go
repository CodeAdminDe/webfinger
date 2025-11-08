//###################################################################//
//#  (c) 2025 Frederic Roggon                                       #//
//#                                                                 #//
//#  Licensed under the terms of GNU AFFERO GENERAL PUBLIC LICENSE. #//
//#  The full terms are provided via LICENSE file which is based    #//
//#  in the root of the code repository.                            #//
//#                                                                 #//
//#  Author: Frederic Roggon <frederic.roggon@codeadmin.de>         #//
//###################################################################//
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	})

	http.HandleFunc("/_healthz", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "200 OK", http.StatusOK)
	})

	http.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		webfingerHandler(w, r, cfg)
	})

	fmt.Println("Webfinger server build <<BUILD>>")
	fmt.Println("Server starting on port 8080...")

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
