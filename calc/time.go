package calc

import (
	"fmt"

	"github.com/rafaljusiak/daily-dashboard/external"
)

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
