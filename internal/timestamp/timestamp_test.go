package timestamp_test

import (
	"strings"
	"testing"
	"time"

	"powerFactors/internal/timestamp"
)

func TestCalculateTimestampsValidSearchRequest(t *testing.T) {
	tests := map[string]struct {
		reqData  timestamp.RequestData
		response string
	}{
		"Hour Period": {
			reqData: timestamp.RequestData{
				Period:          "1h",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  time.Date(2021, time.July, 14, 20, 46, 03, 0, time.UTC),
				SecondTimestamp: time.Date(2021, time.July, 15, 12, 34, 56, 0, time.UTC),
			},
			response: "20210714T210000Z 20210714T220000Z 20210714T230000Z 20210715T000000Z 20210715T010000Z 20210715T020000Z 20210715T030000Z 20210715T040000Z 20210715T050000Z 20210715T060000Z 20210715T070000Z 20210715T080000Z 20210715T090000Z 20210715T100000Z 20210715T110000Z 20210715T120000Z",
		},
		"Day Period": {
			reqData: timestamp.RequestData{
				Period:          "1d",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  time.Date(2021, time.October, 10, 20, 46, 3, 0, time.UTC),
				SecondTimestamp: time.Date(2021, time.November, 15, 12, 34, 56, 0, time.UTC),
			},
			response: "20211010T210000Z 20211011T210000Z 20211012T210000Z 20211013T210000Z 20211014T210000Z 20211015T210000Z 20211016T210000Z 20211017T210000Z 20211018T210000Z 20211019T210000Z 20211020T210000Z 20211021T210000Z 20211022T210000Z 20211023T210000Z 20211024T210000Z 20211025T210000Z 20211026T210000Z 20211027T210000Z 20211028T210000Z 20211029T210000Z 20211030T210000Z 20211031T220000Z 20211101T220000Z 20211102T220000Z 20211103T220000Z 20211104T220000Z 20211105T220000Z 20211106T220000Z 20211107T220000Z 20211108T220000Z 20211109T220000Z 20211110T220000Z 20211111T220000Z 20211112T220000Z 20211113T220000Z 20211114T220000Z",
		},
		"Month Period": {
			reqData: timestamp.RequestData{
				Period:          "1mo",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  time.Date(2021, time.February, 14, 20, 46, 3, 0, time.UTC),
				SecondTimestamp: time.Date(2021, time.November, 15, 12, 34, 56, 0, time.UTC),
			},
			response: "20210228T220000Z 20210331T210000Z 20210430T210000Z 20210531T210000Z 20210630T210000Z 20210731T210000Z 20210831T210000Z 20210930T210000Z 20211031T220000Z",
		},
		"Year Period": {
			reqData: timestamp.RequestData{
				Period:          "1y",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  time.Date(2018, time.February, 14, 20, 46, 3, 0, time.UTC),
				SecondTimestamp: time.Date(2021, time.November, 15, 12, 34, 56, 0, time.UTC),
			},
			response: "20181231T220000Z 20191231T220000Z 20201231T220000Z",
		},
		"Valid period that return no result": {
			reqData: timestamp.RequestData{
				Period:          "1y",
				Timezone:        "Europe/Athens",
				FirstTimestamp:  time.Date(2018, time.February, 14, 20, 46, 3, 0, time.UTC),
				SecondTimestamp: time.Date(2018, time.February, 14, 20, 56, 3, 0, time.UTC),
			},
			response: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := timestamp.CalculateTimestamps(tc.reqData)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if got, want := strings.Join(res, " "), tc.response; got != want {
				t.Errorf("TestCalculateTimestampsValidSearchRequest: expected %v got %v", want, got)
			}
		})
	}
}
