package su_time

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"
)

const (
	DayTime       = 3600 * 24
	WeekTime      = DayTime * 7
	TimeLayout    = "2006-01-02 15:04:05"
	TimeLayoutV2  = "20060102150405"
	TimeLayoutDay = "20060102"
	TimeLayoutV3  = "2006-01-02T15:04:05.305Z"
)

func TimeUnixToDateFile() string {
	return time.Now().Format(TimeLayoutV2)
}
func TimeUnixToDateDayFile() string {
	return time.Now().Format(TimeLayoutDay)
}

// 获取月份的开始和结束时间戳
func GetMonthStartAndEnd(date string) (start, end int64) {
	dateArr := strings.Split(date, "-")
	var myYear, myMonth string
	myYear = dateArr[0]
	myMonth = dateArr[1]
	// 数字月份必须前置补零
	if len(myMonth) == 1 {
		myMonth = "0" + myMonth
	}
	yInt, _ := strconv.Atoi(myYear)
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(TimeLayout, myYear+"-"+myMonth+"-01 00:00:00", loc)
	newMonth := theTime.Month()
	start = time.Date(yInt, newMonth, 1, 0, 0, 0, 0, time.Local).Unix()
	end = time.Date(yInt, newMonth+1, 1, 0, 0, 0, 0, time.Local).Unix()
	return
}
func DateToTime(date string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(TimeLayout, date, loc)
	return theTime.Unix()
}

func DateToTimeByFormat(date string, format string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(format, date, loc)
	tNow := theTime.Unix()

	return tNow
}

// 获取天的开始和结束时间戳
func GetDayStartAndEnd(date string) (start, end int64) {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(TimeLayout, date+" 00:00:00", loc)
	start = theTime.Unix()
	end = theTime.Unix() + DayTime
	return
}

// WeekToDatetime 获得year年第week周的开始时间
func WeekToDatetime(year int, week int) time.Time {
	return FirstWeekOfYear(year).AddDate(0, 0, (week-1)*7)
}

// 获得year年第week周的开始时间戳和结束时间戳
func WeekToStartAndEnd(date string) (start, end int64) {
	dateArr := strings.Split(date, "-")
	year := cast.ToInt(dateArr[0])
	week := cast.ToInt(dateArr[1])
	startTime := WeekToDatetime(year, week)
	endTime := WeekToDatetime(year, week+1)
	start = startTime.Unix()
	end = endTime.Unix()
	return
}

// FirstWeekOfYear 获得第一周的开始时间
func FirstWeekOfYear(year int) time.Time {
	jan1Time := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	jan1WeekStart := jan1Time.AddDate(0, 0, -int(jan1Time.Weekday()+6)%7)
	if y, _ := jan1WeekStart.ISOWeek(); y < year {
		return jan1WeekStart.AddDate(0, 0, 7)
	} else {
		return jan1WeekStart
	}
}

// TimeUnixToDate 获取时间戳的日期
func TimeUnixToDate(timestamp int64) string {
	if timestamp <= 0 {
		return "-"
	}
	return time.Unix(timestamp, 0).Format(TimeLayoutDay)
}

// TimeUnixToDateSecord 获取时间戳的日期 精确到秒
func TimeUnixToDateSecord(timestamp int64) string {
	if timestamp <= 0 {
		return "-"
	}
	if timestamp > 1000000000000 {
		return time.UnixMilli(timestamp).Format(TimeLayout)
	}
	return time.Unix(timestamp, 0).Format(TimeLayout)
}

// ParseTimeLength 获取时间长度
func ParseTimeLength(val int64) string {
	d := int(val / (DayTime))
	val = val - int64(DayTime*d)
	h := int(val / 3600)
	val = val - int64(3600*h)
	m := int(val / 60)
	s := val - int64(60*m)
	return fmt.Sprintf("%d天%d小时%d分%d秒", d, h, m, s)
}

// GetTimeUtcEndTime 获取UTC时区的结束时间
func TimeDateV3ToTimeDate(date string) string {
	t, _ := time.Parse(TimeLayoutV3, date)
	return t.Format(TimeLayout)
}

// 获取当前时间到第n天0点整的时间差，单位秒
func GetTimeDifference(day int) int64 {
	now := time.Now()
	nowTimeStamp := now.Unix()
	nowStr := now.Format("2006-01-02")
	t2, _ := time.ParseInLocation("2006-01-02", nowStr, time.Local)
	t2TimeStamp := t2.AddDate(0, 0, day).Unix()
	return t2TimeStamp - nowTimeStamp
}

func GetDayUtcStartTime() int64 {
	timeStr := time.Now().UTC().Format("2006-01-02")
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t2, _ := time.Parse("2006-01-02", timeStr)
	return t2.AddDate(0, 0, 0).UTC().Unix()
}

func GetTimeUtcStartTime(timeUnix int64) int64 {
	timeStr := time.Unix(timeUnix, 0).UTC().Format("2006-01-02")
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t2, _ := time.Parse("2006-01-02", timeStr)
	return t2.AddDate(0, 0, 0).UTC().Unix()
}

// 美国冬令时
func Special() bool {
	if time.Now().UTC().Month() >= 11 || time.Now().UTC().Month() <= 3 {
		if time.Now().UTC().Month() == 11 {
			if time.Now().UTC().Day() >= 8 {
				return true
			} else {
				return false
			}
		}
		if time.Now().UTC().Month() == 3 {
			if time.Now().UTC().Day() <= 11 {
				return true
			} else {
				return false
			}
		}
		return true
	}
	return false
}

type TimeZoneResponse struct {
	EUS  string `json:"eus"`  //美东时间
	UTC8 string `json:"utc8"` //北京时间
	UTC  string `json:"utc"`  // 标准时间
}

// 获取当前是周几
func GetWeekDay() int32 {
	switch time.Now().Weekday().String() {
	case "Sunday":
		return 7
	case "Monday":
		return 1
	case "Tuesday":
		return 2
	case "Wednesday":
		return 3
	case "Thursday":
		return 4
	case "Friday":
		return 5
	case "Saturday":
		return 6
	default:
		return 0
	}
}

// CurrentTimestamp
// @description 获取当前的时间戳, 单位秒
func CurrentTimestamp() int64 {
	return time.Now().Unix()
}

// Datetime
// @description 当前日期
func Datetime() string {
	return time.Now().Format(TimeLayout)
}

// DateRangeToList
// @description 将日期范围转成日期列表
func DateRangeToList(dateFrom string, dateTo string, format string, newFormat string) []string {
	var dates []string
	if newFormat == "" {
		newFormat = format
	}

	from, _ := time.Parse(format, dateFrom)
	to, _ := time.Parse(format, dateTo)
	for {
		dates = append(dates, from.Format(newFormat))
		from = from.AddDate(0, 0, 1)
		if from.After(to) {
			break
		}

	}

	return dates
}

// TimeRangeToList
// @description 将时间范围转成日期列表
func TimeRangeToDateList(timeFrom int64, timeTo int64, format string) []string {
	var dates []string

	from := time.Unix(timeFrom, 0)
	to := time.Unix(timeTo, 0)
	for {
		dates = append(dates, from.Format(format))
		from = from.AddDate(0, 0, 1)
		if from.After(to) {
			break
		}

	}

	return dates
}

// GetThisWeekStartUnix 获取本周周一的开始时间戳
func GetThisWeekStartUnix() int64 {
	// 获取当前时间
	now := time.Now().UTC()
	// 计算当前时间是周几
	// weekday := int(now.Weekday())
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	// 计算本周的开始时间
	weekStart := now.AddDate(0, 0, offset)

	startUnix := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC).Unix()
	return startUnix
}

// 根据utc+n获取时间
func ConvertTimestampToTime(timestamp int64, offset int, format string) (string, error) {
	loc := time.FixedZone("", offset*60*60)
	t := time.Unix(timestamp, 0).In(loc)
	formattedTime := t.Format(format)
	return formattedTime, nil
}
