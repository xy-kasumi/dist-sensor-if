package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"go.bug.st/serial"
)

type CommandRequest struct {
	Command int `json:"command"`
}

type CommandResponse struct {
	Ok    bool   `json:"ok"`
	Data  int    `json:"data"`
	Error string `json:"error,omitempty"`
}

type ReadRequest struct {
	Addr int `json:"addr"`
}

type ReadResponse struct {
	Ok    bool   `json:"ok"`
	Data  int    `json:"data"`
	Error string `json:"error,omitempty"`
}

type WriteRequest struct {
	Addr int `json:"addr"`
	Data int `json:"data"`
}

type WriteResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(CommandResponse{Ok: false, Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CommandResponse{Ok: true, Data: 0})
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	var req ReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ReadResponse{Ok: false, Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ReadResponse{Ok: true, Data: 0})
}

func handleWrite(w http.ResponseWriter, r *http.Request) {
	var req WriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(WriteResponse{Ok: false, Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(WriteResponse{Ok: true})
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	var filePath string
	switch r.URL.Path {
	case "/":
		filePath = "static/index.html"
	case "/main.js":
		filePath = "static/main.js"
	default:
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

func main() {
	// Flag resolution
	portName := flag.String("port", "/dev/ttyUSB0", "Serial port name")
	baud := flag.Int("baud", 9600, "Serial port baud rate")
	addr := flag.String("addr", ":80", "HTTP listen address")
	flag.Parse()

	mode := &serial.Mode{BaudRate: *baud}
	_, err := serial.Open(*portName, mode)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}

	http.HandleFunc("/api/command", handleCommand)
	http.HandleFunc("/api/read", handleRead)
	http.HandleFunc("/api/write", handleWrite)
	http.HandleFunc("/", handleStatic)
	http.ListenAndServe(*addr, nil)
}
