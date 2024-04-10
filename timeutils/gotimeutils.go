package timeutils

import (
	"fmt"
	"strings"
	"time"
)

func ConvetTimestamsToLocalTime(tm int64) time.Time {
	timeTmp := time.Unix(tm, 0)
	local, _ := time.LoadLocation("Local")
	return timeTmp.In(local)
}

func TimeNowUTC() string {
	//			2021-03-11 01:49:58.968944707 +0000 UTC
	tar := strings.Split(time.Now().UTC().String(), " ")
	return fmt.Sprintf("%s %s", tar[0], tar[1])
}

// Asia/Ho_Chi_Minh
// America/Los_Angeles

func GetTimeStamp(location string) string {
	if location == "" {
		location = "Asia/Ho_Chi_Minh"
	}
	loc, _ := time.LoadLocation(location)
	t := time.Now().In(loc)
	return t.Format("20060102150405")
}

func GetTodaysDate(location string) string {
	if location == "" {
		location = "Asia/Ho_Chi_Minh"
	}
	loc, _ := time.LoadLocation(location)
	current_time := time.Now().In(loc)
	return current_time.Format("2006-01-02")
}

func GetTodaysDateTime(location string) string {
	if location == "" {
		location = "Asia/Ho_Chi_Minh"
	}
	loc, _ := time.LoadLocation(location)
	current_time := time.Now().In(loc)
	return current_time.Format("2006-01-02 15:04:05")
}

func GetTodaysDateTimeFormatted(location string) string {
	if location == "" {
		location = "Asia/Ho_Chi_Minh"
	}
	loc, _ := time.LoadLocation(location)
	current_time := time.Now().In(loc)
	return current_time.Format("Jan 2, 2006 at 3:04 PM")
}

func GetTimeStampFromDate(dtformat string) string {
	form := "Jan 2, 2006 at 3:04 PM"
	t2, _ := time.Parse(form, dtformat)
	return t2.Format("20060102150405")
}
