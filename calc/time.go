package calc

import (
	"fmt"
	"time"

	"github.com/rafaljusiak/daily-dashboard/calc"
	"github.com/rafaljusiak/daily-dashboard/external"
)

func FirstDayOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func SumDuration(timeEntries []external.ClockifyTimeEntryData) (int, error) {
	minutes := 0
	for _, timeEntry := range timeEntries {
		convertedDuration, err := external.ConvertDurationToMinutes(timeEntry.TimeInterval.Duration)
		if err != nil {
			return 0, err
		}
		minutes += convertedDuration
	}
	return minutes, nil
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
	firstDayOfMonth := calc.FirstDayOfMonth(now)
	_ = firstDayOfMonth
	return 0
}
