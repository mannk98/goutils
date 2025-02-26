package timeutils

import (
	"fmt"
	"strings"
	"time"
)

const (
	//value for current date format
	YYYYMMDD = "20060102"

	//value for current datetime format
	YYYYMMDDHHMMSS = "20060102150405"

	//value for current date format
	YYYY_MM_DD = "2006-01-02"

	//value for current datetime format
	YYYY_MM_DD_HH_MM_SS = "2006-01-02-15-04-05"

	Month_day_Year_at_Hour_Minute = "Jan 2, 2006 at 3:04 PM"

	YYYY_MM_DD_HH_MM_SS_WithSpace = "2006-01-02 15:04:05"
)

/*
	type Yyyymmdd struct {
		value string
	}

func NewYyyymmdd(s string) (*Yyyymmdd, error) {

	err := ValidateYmd(s)

	if err != nil {
		return nil, err
	} else {
		return &Yyyymmdd{value: s}, nil
	}

}

	func (item *Yyyymmdd) String() string {
		return item.value
	}

	func (item *Yyyymmdd) Set(s string) error {
		item.value = s
		if !(item.validate(s)) {
			return errors.New("invalid input for Yyyymmdd:" + s)
		} else {
			return nil
		}
	}

	func (item *Yyyymmdd) Scan(arg interface{}) error {
		if arg == nil {
			item.value = ""
		}

		if reflect.TypeOf(arg).String() == "string" {
			s := string((arg.([]byte)))
			if !(item.validate(s)) {
				return errors.New("invalid input for Yyyymmdd:" + s)
			} else {
				item.value = s
			}
		}
		return nil
	}

	func (item *Yyyymmdd) UnmarshalJSON(b []byte) error {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if !(item.validate(s)) {
			return errors.New("invalid input for Yyyymmdd:" + s)
		} else {
			item.value = s
		}
		return nil
	}

	func (item *Yyyymmdd) MarshalJSON() ([]byte, error) {
		return json.Marshal(item.value)
	}

	func (item *Yyyymmdd) validate(s string) bool {
		err := ValidateYmd(s)
		if err != nil {
			return false
		} else {
			return true
		}
	}

	func ValidateYmd(s string) error {
		if len(s) != 8 {
			return errors.New("YYYYMMDD length must be 8")
		}

		_, err := strconv.Atoi(s)
		if err != nil {
			return errors.New("YYYYMMDD must be numbers")
		}

		reg := regexp.MustCompile(`[-|/|:| |　]`)
		str := reg.ReplaceAllString(s, "")
		format := string([]rune("20060102150405")[:len(str)])
		_, errParse := time.Parse(format, str)

		if errParse != nil {
			return errors.New("YYYYMMDD must be valid comnination with year, month and day:" + errParse.Error())
		}

		return nil
	}
*/

/*
timeStamp: unixTime
*/
func GetLocalTimeFromTimestamp(timeStamp int64) string {
	timeTmp := time.Unix(timeStamp, 0)
	local, _ := time.LoadLocation("Local")
	return timeTmp.In(local).String()
}

/*
default: Asia/Ho_Chi_Minh GMT +7
*/
func GetTimeNow_YYYYMMDDHHMMSS(location string) string {
	if location == "" {
		location = "Asia/Ho_Chi_Minh"
	}
	loc, _ := time.LoadLocation(location)
	t := time.Now().In(loc)
	return t.Format(YYYYMMDDHHMMSS)
}

func GetTimeNow_YYYY_MM_DD(location string) string {
	if location == "" {
		location = "Asia/Ho_Chi_Minh"
	}
	loc, _ := time.LoadLocation(location)
	current_time := time.Now().In(loc)
	return current_time.Format(YYYY_MM_DD)
}

func GetTimeNow_YYYY_MM_DD_HH_MM_SS_WithSpace(location string) string {
	if location == "" {
		location = "Asia/Ho_Chi_Minh"
	}
	loc, _ := time.LoadLocation(location)
	current_time := time.Now().In(loc)
	return current_time.Format(YYYY_MM_DD_HH_MM_SS_WithSpace)
}

/*
Format: 2009-11-10 23:00:00
*/
func GetTimeNowUTC_YYYY_MM_DD_HH_MM_SS_WithSpace() string {
	tar := strings.Split(time.Now().UTC().String(), " ")
	return fmt.Sprintf("%s %s", tar[0], tar[1])
}

func GetTimeNow_Month_day_Year_at_Hour_Minute(location string) string {
	if location == "" {
		location = "Asia/Ho_Chi_Minh"
	}
	loc, _ := time.LoadLocation(location)
	current_time := time.Now().In(loc)
	return current_time.Format(Month_day_Year_at_Hour_Minute)
}

/*
GetTimeStampFromDate("Feb 25, 2025 at 10:30 AM")
*/
func Get_YYYYMMDDHHMMSS_From_Month_day_Year_at_Hour_Minute(dtformat string) string {
	form := Month_day_Year_at_Hour_Minute
	t2, _ := time.Parse(form, dtformat)
	return t2.Format(YYYYMMDDHHMMSS)
}
