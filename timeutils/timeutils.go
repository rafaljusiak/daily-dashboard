package timeutils

import (
	"fmt"
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
	return fmt.Sprintf("%dh %02dm", h, m)
}

func MinutesToHours(minutes int) float64 {
	return float64(minutes) / 60.0
}

func WorkingHoursForCurrentMonth() int {
	now := time.Now()
	firstDayOfMonth := FirstDayOfMonth(now)
	_ = firstDayOfMonth
	return 0
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
