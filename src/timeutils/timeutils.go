package timeutils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

func FirstDayOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func MinutesToString(minutes int) string {
	h := minutes / 60
	m := minutes % 60
	mAbs := int(math.Abs(float64(m)))
	return fmt.Sprintf("%dh %02dm", h, mAbs)
}

func HoursToString(hours int) string {
	return fmt.Sprintf("%dh 00m", hours)
}

func MinutesToHours(minutes int) float64 {
	return float64(minutes) / 60.0
}

func WorkingHoursForCurrentMonth() int {
	now := time.Now()
	firstDayOfMonth := FirstDayOfMonth(now)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	return WorkingHoursBetweenDates(firstDayOfMonth, lastDayOfMonth)
}

func WorkingHoursUntilToday() int {
	now := time.Now()
	firstDayOfMonth := FirstDayOfMonth(now)

	return WorkingHoursBetweenDates(firstDayOfMonth, now)
}

func WorkingHoursBetweenDates(from, to time.Time) int {
	counter := 0
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		if d.Weekday()%6 != 0 {
			counter++
		}
	}

	return counter * 8
}

func ConvertDurationToMinutes(duration string) (int, error) {
	if len(duration) == 0 {
		return 0, nil
	}

	re := regexp.MustCompile(`PT((?P<hours>\d+)H)?((?P<minutes>\d+)M)?`)
	matches := re.FindStringSubmatch(duration)
	if matches == nil {
		return 0, nil
	}

	hours, _ := strconv.Atoi(matches[re.SubexpIndex("hours")])
	minutes, _ := strconv.Atoi(matches[re.SubexpIndex("minutes")])

	return hours*60 + minutes, nil
}
