package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func getParam(r *http.Request, s string) string {
	keys, ok := r.URL.Query()[s]
	if !ok || len(keys) < 1 {
		log.Printf("Url Param %v is missing", s)
		return ""
	}
	return keys[0]
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func sendUnathorized(w http.ResponseWriter) {
	var resp errResp
	resp.Error = "unauthorized"
	respW, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respW)
}

func sendMessage(w http.ResponseWriter, s string) {
	var resp succResp
	resp.Message = s
	respW, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respW)
}
