package main

import (
	"time"

	"github.com/vjeantet/jodaTime"
)

const year = "YYYY"
const yearMonth = "YYYY.MM"
const yearMonthDay = "YYYY.MM.dd"
const yearMonthDayHour = "YYYY.MM.ddTHHZ"
const yearMonthDayHourMinute = "YYYY.MM.ddTHH:mmZ"
const yearMonthDayHourMinuteSecond = "YYYY.MM.ddTHH:mm:ssZ"

func getEpochNow() int64 {
	return time.Now().Unix()
}

func getToday() string {
	return jodaTime.Format("YYYY.MM.dd", time.Now())
}

//time getters
func getYear() string {
	return jodaTime.Format(year, time.Now())
}
func getMonth() string {
	return jodaTime.Format(yearMonth, time.Now())
}
func getDay() string {
	return jodaTime.Format(yearMonthDay, time.Now())
}
func getHour() string {
	return jodaTime.Format(yearMonthDayHour, time.Now())
}
func getMinute() string {
	return jodaTime.Format(yearMonthDayHourMinute, time.Now())
}
func getSecond() string {
	return jodaTime.Format(yearMonthDayHourMinuteSecond, time.Now())
}

//time past/future getters for year,month,day
func getTimeXYearsAgo(x int, format string) string {
	return jodaTime.Format(format, time.Now().AddDate(-x, 0, 0))
}
func getTimeXMonthsAgo(x int, format string) string {
	return jodaTime.Format(format, time.Now().AddDate(0, -x, 0))
}
func getTimeXDaysAgo(x int, format string) string {
	return jodaTime.Format(format, time.Now().AddDate(0, 0, -x))
}
func getTimeInXYears(x int, format string) string {
	return jodaTime.Format(format, time.Now().AddDate(x, 0, 0))
}
func getTimeInXMonths(x int, format string) string {
	return jodaTime.Format(format, time.Now().AddDate(0, x, 0))
}
func getTimeInXDays(x int, format string) string {
	return jodaTime.Format(format, time.Now().AddDate(0, 0, x))
}

//time past/future getters for hour,minute,second
func getTimeXHoursAgo(x int, format string) string {
	return jodaTime.Format(format, time.Now().Add(time.Duration(-x)*time.Hour))
}
func getTimeXMinAgo(x int, format string) string {
	return jodaTime.Format(format, time.Now().Add(time.Duration(-x)*time.Minute))
}
func getTimeXSecondsAgo(x int, format string) string {
	return jodaTime.Format(format, time.Now().Add(time.Duration(-x)*time.Second))
}
func getTimeInXHours(x int, format string) string {
	return jodaTime.Format(format, time.Now().Add(time.Duration(x)*time.Hour))
}
func getTimeInXMin(x int, format string) string {
	return jodaTime.Format(format, time.Now().Add(time.Duration(x)*time.Minute))
}
func getTimeInXSeconds(x int, format string) string {
	return jodaTime.Format(format, time.Now().Add(time.Duration(x)*time.Second))
}
