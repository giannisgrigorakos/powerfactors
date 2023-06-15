package timestamp

import (
	"fmt"
	"time"
)

const (
	ISO8601Format = "20060102T150405Z"
)

// RequestData exist to separate the business layer from the server layer.
// In this struct we kep only the necessary information for the business layer.
type RequestData struct {
	Period          string
	Timezone        string
	FirstTimestamp  time.Time
	SecondTimestamp time.Time
}

func addTime(prd string, t *time.Time) time.Time {
	switch prd {
	case "1h":
		*t = t.Add(time.Hour)
	case "1d":
		*t = t.AddDate(0, 0, 1)
	case "1mo":
		// We want to advance to the last day of the next month. So we advance two months
		// and, we go back one day to achieve this functionality.
		*t = t.AddDate(0, 1, 0)
	case "1y":
		*t = t.AddDate(1, 0, 0)
	default:
	}
	return *t
}

// CalculateTimestamps calculates a series of timestamps based on the given parameters.
// It generates timestamps within the specified time range and time period, adjusted to the provided time zone.
//
// Parameters:
//
//	reqData.FirstTimestamp       : The starting timestamp.
//	t2       : The ending timestamp.
//	period   : The time period for generating timestamps. Valid values are "1h" (hourly),
//	           "1d" (daily), "1mo" (monthly), and "1y" (yearly).
//	timezone : The time zone used for the calculations.
//
// Returns:
//
//	A slice of strings containing the generated timestamps.
//	An error if there was an issue with the time zone loading.
//
// The idea here is that we get the time in UTC, we convert it to the local time with the timezone
// that is provided and then before writing to the result we convert it back to UTC.
func CalculateTimestamps(reqData RequestData) ([]string, error) {
	location, err := time.LoadLocation(reqData.Timezone)
	if err != nil {
		return []string{}, fmt.Errorf("error loading time zone: %w", err)
	}

	switch reqData.Period {
	case "1h":
		reqData.FirstTimestamp = reqData.FirstTimestamp.Truncate(time.Hour)
		addTime(reqData.Period, &reqData.FirstTimestamp)
	case "1d":
		reqData.FirstTimestamp = reqData.FirstTimestamp.Truncate(time.Hour).Add(time.Hour)
	case "1mo":
		year, month, _ := reqData.FirstTimestamp.AddDate(0, 1, 0).Date()
		reqData.FirstTimestamp = time.Date(reqData.FirstTimestamp.Year(), reqData.FirstTimestamp.Month(), reqData.FirstTimestamp.Day(), reqData.FirstTimestamp.Hour(), reqData.FirstTimestamp.Minute(), reqData.FirstTimestamp.Second(), 0, location)
		reqData.FirstTimestamp = time.Date(year, month, 1, 0, 0, 0, 0, reqData.FirstTimestamp.Location())
	case "1y":
		reqData.FirstTimestamp = time.Date(reqData.FirstTimestamp.Year()+1, 1, 1, 0, 0, 0, 0, location)
	}

	// T1 and T2 should be converted to the timezone provided in the request.
	reqData.FirstTimestamp = reqData.FirstTimestamp.In(location)
	reqData.SecondTimestamp = reqData.SecondTimestamp.In(location)

	result := []string{}
	for {
		if reqData.FirstTimestamp.After(reqData.SecondTimestamp) {
			break
		}
		// Before we append the result we convert the timestamp back to UTC.
		result = append(result, reqData.FirstTimestamp.In(time.UTC).Format(ISO8601Format))
		addTime(reqData.Period, &reqData.FirstTimestamp)
	}
	return result, nil
}
