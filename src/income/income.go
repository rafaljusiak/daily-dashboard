package income

func Calculate(timeHours float64, hourlyRate float64, exchangeRate float64) float64 {
	return timeHours * hourlyRate * exchangeRate
}
