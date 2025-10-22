package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
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

type Server struct {
	dev *CD22
}

func registerJsonHandler[ReqT any, RespT any](path string, exec func(*ReqT) (*RespT, error)) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Decode & validate
		var req ReqT
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "invalid JSON: %v", err)
			return
		}

		// Execute
		resp, err := exec(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send response as JSON
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}

func (s *Server) execCommand(req *CommandRequest) (*CommandResponse, error) {
	data, err := s.dev.Send(CMD_COMMAND, uint16(req.Command))
	if err != nil {
		return &CommandResponse{Ok: false, Error: err.Error()}, nil
	}
	return &CommandResponse{Ok: true, Data: int(data)}, nil
}

func (s *Server) execRead(req *ReadRequest) (*ReadResponse, error) {
	data, err := s.dev.Send(CMD_READ, uint16(req.Addr))
	if err != nil {
		return &ReadResponse{Ok: false, Error: err.Error()}, nil
	}
	return &ReadResponse{Ok: true, Data: int(data)}, nil
}

func (s *Server) execWrite(req *WriteRequest) (*WriteResponse, error) {
	// Set address by reading
	_, err := s.dev.Send(CMD_READ, uint16(req.Addr))
	if err != nil {
		return &WriteResponse{Ok: false, Error: err.Error()}, nil
	}
	// Write
	_, err = s.dev.Send(CMD_WRITE, uint16(req.Data))
	if err != nil {
		return &WriteResponse{Ok: false, Error: err.Error()}, nil
	}
	return &WriteResponse{Ok: true}, nil
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
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
	portName := flag.String("port", "/dev/ttyUSB0", "Serial port name")
	baud := flag.Int("baud", 9600, "Serial port baud rate")
	addr := flag.String("addr", ":80", "HTTP listen address")
	flag.Parse()

	dev, err := NewCD22(*portName, *baud)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}

	server := &Server{dev: dev}
	registerJsonHandler("/api/command", server.execCommand)
	registerJsonHandler("/api/read", server.execRead)
	registerJsonHandler("/api/write", server.execWrite)
	http.HandleFunc("/", server.handleStatic)
	log.Printf("Listening on %s (serial: %s @ %d baud)", *addr, *portName, *baud)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
