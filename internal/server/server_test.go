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
	if got, want := string(data), "{\"status\":\"error\",\"desc\":\"Bad Request:EOF\"}\n"; got != want {
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
		"Day Period": {
			request: server.SearchRequest{
				Period:          "1d",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20211010T204603Z",
				SecondTimestamp: "20211115T123456Z",
			},
			response: "{\"Timestamps\":[\"20211010T210000Z\",\"20211011T210000Z\",\"20211012T210000Z\",\"20211013T210000Z\",\"20211014T210000Z\",\"20211015T210000Z\",\"20211016T210000Z\",\"20211017T210000Z\",\"20211018T210000Z\",\"20211019T210000Z\",\"20211020T210000Z\",\"20211021T210000Z\",\"20211022T210000Z\",\"20211023T210000Z\",\"20211024T210000Z\",\"20211025T210000Z\",\"20211026T210000Z\",\"20211027T210000Z\",\"20211028T210000Z\",\"20211029T210000Z\",\"20211030T210000Z\",\"20211031T220000Z\",\"20211101T220000Z\",\"20211102T220000Z\",\"20211103T220000Z\",\"20211104T220000Z\",\"20211105T220000Z\",\"20211106T220000Z\",\"20211107T220000Z\",\"20211108T220000Z\",\"20211109T220000Z\",\"20211110T220000Z\",\"20211111T220000Z\",\"20211112T220000Z\",\"20211113T220000Z\",\"20211114T220000Z\"]}\n",
		},
		"Month Period": {
			request: server.SearchRequest{
				Period:          "1mo",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20210214T204603Z",
				SecondTimestamp: "20211115T123456Z",
			},
			response: "{\"Timestamps\":[\"20210228T220000Z\",\"20210331T210000Z\",\"20210430T210000Z\",\"20210531T210000Z\",\"20210630T210000Z\",\"20210731T210000Z\",\"20210831T210000Z\",\"20210930T210000Z\",\"20211031T220000Z\"]}\n",
		},
		"Year Period": {
			request: server.SearchRequest{
				Period:          "1y",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20180214T204603Z",
				SecondTimestamp: "20211115T123456Z",
			},
			response: "{\"Timestamps\":[\"20181231T220000Z\",\"20191231T220000Z\",\"20201231T220000Z\"]}\n",
		},
		"Valid period that return no result": {
			request: server.SearchRequest{
				Period:          "1y",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  "20180214T204603Z",
				SecondTimestamp: "20180214T205003Z",
			},
			response: "{}\n",
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
