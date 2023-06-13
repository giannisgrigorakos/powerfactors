package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	iso8601Format = "20060102T150405Z"
)

type Server struct {
	Router *mux.Router
}

type SearchRequest struct {
	Period          string `json:"period"`
	Timezone        string `json:"tz"`
	FirstTimestamp  string `json:"t1"`
	SecondTimestamp string `json:"t2"`
}

// should I break it in failresponce and successresponse or just have it as one?
type SearchResponse struct {
	Timestamps  []string `json:",omitempty"`
	Status      string   `json:"status,omitempty"`
	Description string   `json:"desc,omitempty"`
}

func addTime(prd string, t *time.Time) time.Time {
	switch prd {
	case "1h":
		*t = t.Add(time.Hour)
		return *t
	case "1d":
		*t = t.AddDate(0, 0, 1)
		return *t
	case "1mo":
		*t = t.AddDate(0, 2, -t.Day()).Truncate(time.Hour * 24)
		return *t
	case "1y":
		*t = t.AddDate(1, 0, 0)
		return *t
	default:
		panic("Should not be this case ever as validation must have caught this") //TODO: check this
	}
}

var supportedPeriods map[string]string

// TODO: documentation
func NewServer() *Server {
	// Although we won't use the values of this map we use it because we want to take advantage of the constant lookup time.
	supportedPeriods = map[string]string{
		"1h":  "1 hour",
		"1d":  "1 day",
		"1mo": "1 month",
		"1y":  "1 year",
	}

	server := &Server{
		Router: mux.NewRouter(),
	}

	server.Router.HandleFunc("/ptlist", server.FindMatchingTimestamps).Methods("GET")

	return server
}

// TODO: Documentation
func (s *Server) FindMatchingTimestamps(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	// TODO: why should I do this?
	defer r.Body.Close()
	if err != nil {
		// for the sake of this project we are just going to print the error and return 400 to the client.
		// Otherwise we would use a logger in order log this error
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(w)
		encoder.Encode(&SearchResponse{
			Status:      "error",
			Description: "Bad Request:" + err.Error(),
		})
		return
	}

	t1, t2, err := validateRequestAndFetchTimestamps(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(w)
		encoder.Encode(&SearchResponse{
			Status:      "error",
			Description: err.Error(),
		})
		return

	}

	seconds, minutes, hours, days, months, years := t1.Second(), t1.Minute(), t1.Hour(), t1.Day(), t1.Month(), t1.Year()

	switch req.Period {
	case "1h":
		if seconds != 0 || minutes != 0 {
			// we set to zero the seconds and minutes, and we advance the hour to the next hour
			t1 = t1.Truncate(time.Hour)
			addTime(req.Period, &t1)
		}
	case "1d":
		if seconds != 0 || minutes != 0 || hours != 0 {
			// we set to zero the seconds, minutes and hours, and we advance the day to the next day
			t1 = t1.Truncate(time.Hour * 24)
			addTime(req.Period, &t1)
		}
	case "1mo":
		if seconds != 0 || minutes != 0 || hours != 0 || days != 0 {
			t1 = t1.AddDate(0, 1, -t1.Day()).Truncate(time.Hour * 24)
		}
	case "1y":
		if seconds != 0 || minutes != 0 || hours != 0 || days != 0 || months != 0 {
			t1 = time.Date(years, 12, 31, 0, 0, 0, 0, time.UTC)
		}
	}

	// we do it like this, so we have empty array and not nil
	result := []string{}
	for {
		if t1.After(t2) {
			break
		}
		result = append(result, t1.Format(iso8601Format))
		addTime(req.Period, &t1)
	}

	fmt.Println(result)

	encoder := json.NewEncoder(w)
	encoder.Encode(&SearchResponse{
		Timestamps: result,
	})
}

func validateRequestAndFetchTimestamps(req SearchRequest) (time.Time, time.Time, error) {
	// Check if fields are empty and if they are that means they are missing from the request.
	if req.Period == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("period field `period` is missing from the request: %v", req)
	}
	if req.Timezone == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("timezone field `tz` is missing from the request: %v", req)
	}
	if req.FirstTimestamp == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("first timestamp field `t1` is missing from the request: %v", req)
	}
	if req.SecondTimestamp == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("second timestamp field `t2` is missing from the request: %v", req)
	}

	// Validate that the specific period is supported by the system.
	if _, ok := supportedPeriods[req.Period]; !ok {
		return time.Time{}, time.Time{}, fmt.Errorf("unsupported period: %v", req.Period)
	}

	t1, err := time.Parse(iso8601Format, req.FirstTimestamp)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("could not parse t1: %v as it's not in 20060102T150405Z (ISO8601) format", req.FirstTimestamp)
	}
	t2, err := time.Parse(iso8601Format, req.SecondTimestamp)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("could not parse t2: %v as it's not in 20060102T150405Z (ISO8601) format", req.SecondTimestamp)
	}
	if t1.After(t2) {
		return time.Time{}, time.Time{}, fmt.Errorf("t1: %v should be before t2: %v", req.FirstTimestamp, req.SecondTimestamp)
	}
	return t1, t2, nil
}
