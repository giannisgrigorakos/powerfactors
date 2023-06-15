package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"powerFactors/internal/server"
)

func TestFindMatchingTimestampsEmptyBodySearchRequest(t *testing.T) {
	server := &server.Server{
		Router: mux.NewRouter(),
	}

	responseRecorder := httptest.NewRecorder()
	server.FindMatchingTimestamps(responseRecorder, httptest.NewRequest("GET", "/ptlist", nil))
	res := responseRecorder.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if got, want := string(data), "{\"status\":\"error\",\"desc\":\"Bad Request/ Empty Request:EOF\"}\n"; got != want {
		t.Errorf("TestFindMatchingTimestampsEmptySearchRequest Data: expected %v got %v", want, got)
	}
	if got, want := res.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("TestFindMatchingTimestampsEmptySearchRequest Status Code: expected %v got %v", want, got)
	}
}

func TestFindMatchingTimestampsInvalidSearchRequest(t *testing.T) {
	srv := server.NewServer()
	tests := map[string]struct {
		wrongRequest server.SearchRequest
		expectedErr  string
	}{
		"Missing Period": {
			wrongRequest: server.SearchRequest{
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20180214T204603Z",
				SecondTimestamp: "20180214T204703Z",
			},
			expectedErr: "{\"status\":\"error\",\"desc\":\"period field `period` is missing from the request: { Europe/Athens 20180214T204603Z 20180214T204703Z}\"}\n",
		},
		"Missing Timezone": {
			wrongRequest: server.SearchRequest{
				Period:          "1y",
				FirstTimestamp:  "20180214T204603Z",
				SecondTimestamp: "20180214T204703Z",
			},
			expectedErr: "{\"status\":\"error\",\"desc\":\"timezone field `tz` is missing from the request: {1y  20180214T204603Z 20180214T204703Z}\"}\n",
		},
		"Missing First Timestamp": {
			wrongRequest: server.SearchRequest{
				Period:          "1y",
				Timezone:        "Europe/Athens",
				SecondTimestamp: "20180214T204703Z",
			},
			expectedErr: "{\"status\":\"error\",\"desc\":\"first timestamp field `t1` is missing from the request: {1y Europe/Athens  20180214T204703Z}\"}\n",
		},
		"Missing Second Timestamp": {
			wrongRequest: server.SearchRequest{
				Period:         "1y",
				Timezone:       "Europe/Athens",
				FirstTimestamp: "20180214T204603Z",
			},
			expectedErr: "{\"status\":\"error\",\"desc\":\"second timestamp field `t2` is missing from the request: {1y Europe/Athens 20180214T204603Z }\"}\n",
		},
		"Unsupported Period": {
			wrongRequest: server.SearchRequest{
				Period:          "1ya",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20180214T204603Z",
				SecondTimestamp: "20180214T204703Z",
			},
			expectedErr: "{\"status\":\"error\",\"desc\":\"unsupported period: 1ya\"}\n",
		},
		"Unsupported Timezone": {
			wrongRequest: server.SearchRequest{
				Period:          "1y",
				Timezone:        "Eurasia",
				FirstTimestamp:  "20180214T204603Z",
				SecondTimestamp: "20180214T204703Z",
			},
			expectedErr: "{\"status\":\"error\",\"desc\":\"error loading time zone: unknown time zone Eurasia\"}\n",
		},
		"First timestamp wrong format": {
			wrongRequest: server.SearchRequest{
				Period:          "1y",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20180214T20460asd",
				SecondTimestamp: "20180214T204703Z",
			},
			expectedErr: "{\"status\":\"error\",\"desc\":\"could not parse t1: 20180214T20460asd as it's not in 20060102T150405Z (ISO8601) format\"}\n",
		},
		"Second timestamp wrong format": {
			wrongRequest: server.SearchRequest{
				Period:          "1y",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20180214T204603Z",
				SecondTimestamp: "20180214T204asd",
			},
			expectedErr: "{\"status\":\"error\",\"desc\":\"could not parse t2: 20180214T204asd as it's not in 20060102T150405Z (ISO8601) format\"}\n",
		},
		"First timestamp is after second timestamp": {
			wrongRequest: server.SearchRequest{
				Period:          "1y",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20180214T204603Z",
				SecondTimestamp: "20180214T204503Z",
			},
			expectedErr: "{\"status\":\"error\",\"desc\":\"t1: 20180214T204603Z should be before t2: 20180214T204503Z\"}\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			body, _ := json.Marshal(tc.wrongRequest)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/ptlist", bytes.NewReader(body))
			srv.FindMatchingTimestamps(w, r)
			res := w.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if got, want := string(data), tc.expectedErr; got != want {
				t.Errorf("TestFindMatchingTimestampsInvalidSearchRequest Data: expected %v got %v", want, got)
			}
			if got, want := res.StatusCode, http.StatusBadRequest; got != want {
				t.Errorf("TestFindMatchingTimestampsInvalidSearchRequest Status Code: expected %v got %v", want, got)
			}
		})
	}
}

func TestFindMatchingTimestampsValidSearchRequest(t *testing.T) {
	srv := server.NewServer()
	tests := map[string]struct {
		request  server.SearchRequest
		response string
	}{
		"Hour Period": {
			request: server.SearchRequest{
				Period:          "1h",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20210714T204603Z",
				SecondTimestamp: "20210715T123456Z",
			},
			response: "{\"Timestamps\":[\"20210714T210000Z\",\"20210714T220000Z\",\"20210714T230000Z\",\"20210715T000000Z\",\"20210715T010000Z\",\"20210715T020000Z\",\"20210715T030000Z\",\"20210715T040000Z\",\"20210715T050000Z\",\"20210715T060000Z\",\"20210715T070000Z\",\"20210715T080000Z\",\"20210715T090000Z\",\"20210715T100000Z\",\"20210715T110000Z\",\"20210715T120000Z\"]}\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			body, _ := json.Marshal(tc.request)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/ptlist", bytes.NewReader(body))
			srv.FindMatchingTimestamps(w, r)
			res := w.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if got, want := string(data), tc.response; got != want {
				t.Errorf("TestFindMatchingTimestampsValidSearchRequest: expected %v got %v", want, got)
			}
			if got, want := res.StatusCode, http.StatusOK; got != want {
				t.Errorf("TestFindMatchingTimestampsValidSearchRequest Status Code: expected %v got %v", want, got)
			}
		})
	}
}
