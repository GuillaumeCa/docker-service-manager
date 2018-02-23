package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getEnv(env, std string) string {
	e := os.Getenv(env)
	if e == "" {
		return std
	}
	return e
}

func resJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("could not encode payload: %v", err)
		return
	}
	w.Write(body)
}

func resStatus(w http.ResponseWriter, status int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.WriteHeader(status)
}
