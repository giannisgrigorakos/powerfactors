package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"powerFactors/internal/timestamp"
)

// Server contains the Router that will be used. It can be extended with adding extra fields that
// based on the specific features that we want to develop. For example if we want to connect with a
// database we would create another layer and then the handler would call the server.SomeManager.DoSomethingInTheDB.
type Server struct {
	Router *mux.Router
}

// SearchRequest holds all the necessary information for a successful request from the clients.
type SearchRequest struct {
	Period          string `json:"period"`
	Timezone        string `json:"tz"`
	FirstTimestamp  string `json:"t1"`
	SecondTimestamp string `json:"t2"`
}

// SearchResponse contains all the necessary information for a response to the clients.
// In a successful response only the Timestamps field will be filled and the other two will
// be omitted. In an unsuccessful response the Status and Description fields will be filled
// providing information on what went wrong, to the clients, while the remaining field will
// be omitted.
type SearchResponse struct {
	Timestamps  []string `json:",omitempty"`
	Status      string   `json:"status,omitempty"`
	Description string   `json:"desc,omitempty"`
}

var supportedPeriods map[string]struct{}

// NewServer initializes a new server and instantiates the supported periods.
func NewServer() *Server {
	// Although we won't use the values of this map we use it because we want to take advantage of the constant lookup time.
	// In order to make it more extensible we can add more supported periods as we see fit.
	supportedPeriods = map[string]struct{}{
		"1h":  {},
		"1d":  {},
		"1mo": {},
		"1y":  {},
	}

	server := &Server{
		Router: mux.NewRouter(),
	}

	server.Router.HandleFunc("/ptlist", server.FindMatchingTimestamps).Methods("GET")

	return server
}

// FindMatchingTimestamps handles the HTTP request to find matching timestamps.
// It validates the request, calculates the timestamps based on the provided parameters,
// and returns the result in the HTTP response.
//
// Method:   POST
// Endpoint: /matching-timestamps
//
// Request Body:
//
//	The request body should contain a JSON object with the following structure:
//	{
//	    "t1": "yyyy-MM-ddTHH:mm:ssZ",
//	    "t2": "yyyy-MM-ddTHH:mm:ssZ",
//	    "period": "1h",
//	    "timezone": "UTC"
//	}
//
// Response:
//
//	The response will be a JSON object with the following structure:
//	{
//	    "timestamps": ["2021-10-10T21:00:00Z", "2021-10-11T21:00:00Z", ...],
//	}
//	In case of an error, the response will have the following structure:
//	{
//	    "status": "error",
//	    "description": "Error message"
//	}
func (s *Server) FindMatchingTimestamps(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(w)
		encoder.Encode(&SearchResponse{
			Status:      "error",
			Description: "Bad Request:" + err.Error(),
		})
		return
	}

	t1, t2, period, timezone, err := validateRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(w)
		encoder.Encode(&SearchResponse{
			Status:      "error",
			Description: err.Error(),
		})
		return

	}

	result, err := timestamp.CalculateTimestamps(t1, t2, period, timezone)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(w)
		encoder.Encode(&SearchResponse{
			Status:      "error",
			Description: err.Error(),
		})
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(&SearchResponse{
		Timestamps: result,
	})
}

func validateRequest(req SearchRequest) (time.Time, time.Time, string, string, error) {
	// Check if fields are empty and if they are that means they are missing from the request.
	if req.Period == "" {
		return time.Time{}, time.Time{}, "", "", fmt.Errorf("period field `period` is missing from the request: %v", req)
	}
	if req.Timezone == "" {
		return time.Time{}, time.Time{}, "", "", fmt.Errorf("timezone field `tz` is missing from the request: %v", req)
	}
	if req.FirstTimestamp == "" {
		return time.Time{}, time.Time{}, "", "", fmt.Errorf("first timestamp field `t1` is missing from the request: %v", req)
	}
	if req.SecondTimestamp == "" {
		return time.Time{}, time.Time{}, "", "", fmt.Errorf("second timestamp field `t2` is missing from the request: %v", req)
	}

	// Validate that the specific period is supported.
	if _, ok := supportedPeriods[req.Period]; !ok {
		return time.Time{}, time.Time{}, "", "", fmt.Errorf("unsupported period: %v", req.Period)
	}

	t1, err := time.Parse(timestamp.ISO8601Format, req.FirstTimestamp)
	if err != nil {
		return time.Time{}, time.Time{}, "", "", fmt.Errorf("could not parse t1: %v as it's not in 20060102T150405Z (ISO8601) format", req.FirstTimestamp)
	}
	t2, err := time.Parse(timestamp.ISO8601Format, req.SecondTimestamp)
	if err != nil {
		return time.Time{}, time.Time{}, "", "", fmt.Errorf("could not parse t2: %v as it's not in 20060102T150405Z (ISO8601) format", req.SecondTimestamp)
	}
	if t1.After(t2) {
		return time.Time{}, time.Time{}, "", "", fmt.Errorf("t1: %v should be before t2: %v", req.FirstTimestamp, req.SecondTimestamp)
	}
	return t1, t2, req.Period, req.Timezone, nil
}
