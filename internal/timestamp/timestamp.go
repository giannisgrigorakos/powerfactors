package timestamp

import (
	"fmt"
	"time"
)

const (
	ISO8601Format = "20060102T150405Z"
)

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
//	t1       : The starting timestamp.
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
func CalculateTimestamps(t1, t2 time.Time, period, timezone string) ([]string, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return []string{}, fmt.Errorf("error loading time zone: %w", err)
	}

	switch period {
	case "1h":
		// we set to zero the seconds and minutes, and we advance the hour to the next hour
		t1 = t1.Truncate(time.Hour)
		addTime(period, &t1)
	case "1d":
		// we set to zero the seconds, minutes, and we advance the day to the next hour
		t1 = t1.Truncate(time.Hour).Add(time.Hour)
	case "1mo":
		//
		year, month, _ := t1.AddDate(0, 1, 0).Date()
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, location)
		t1 = time.Date(year, month, 1, 0, 0, 0, 0, t1.Location())
	case "1y":
		//
		t1 = time.Date(t1.Year()+1, 1, 1, 0, 0, 0, 0, location)
	}

	// T1 and T2 should be converted to the timezone provided in the request.
	t1 = t1.In(location)
	t2 = t2.In(location)

	result := []string{}
	for {
		if t1.After(t2) {
			break
		}
		// Before we append the result we convert the timestamp back to UTC.
		result = append(result, t1.In(time.UTC).Format(ISO8601Format))
		addTime(period, &t1)
	}
	return result, nil
}
