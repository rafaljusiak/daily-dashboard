package calc

import (
	"github.com/rafaljusiak/daily-dashboard/external"
	"github.com/rafaljusiak/daily-dashboard/timeutils"
)

func SumDuration(timeEntries []external.ClockifyTimeEntryData) (int, error) {
	minutes := 0
	for _, timeEntry := range timeEntries {
		convertedDuration, err := timeutils.ConvertDurationToMinutes(timeEntry.TimeInterval.Duration)
		if err != nil {
			return 0, err
		}
		minutes += convertedDuration
	}
	return minutes, nil
}
