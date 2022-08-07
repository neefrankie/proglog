package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
)

var decoder = schema.NewDecoder()

func NewHttpServer(addr string) *http.Server {
	httpSrv := newHttpServer()
	r := mux.NewRouter()

	r.HandleFunc("/", httpSrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpSrv.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type HttpServer struct {
	Log *Log
}

func newHttpServer() *HttpServer {
	return &HttpServer{
		Log: NewLog(),
	}
}

// handleProduce creates a record.
// Test data:
// {
//    "record": {
//        "value": "TGV0J3MgR28gIzEK"
//    }
// },
// {
//    "record": {
//        "value": "TGV0J3MgR28gIzIK"
//    }
// },
// {
//    "record": {
//        "value": "TGV0J3MgR28gIzMK"
//    }
// }
func (s *HttpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HttpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req ConsumeRequest
	err = decoder.Decode(&req, r.Form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Record: %v\n", record)

	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset" schema:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}
