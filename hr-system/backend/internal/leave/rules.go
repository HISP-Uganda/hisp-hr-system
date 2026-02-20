package leave

import "time"

func ComputeWorkingDays(startDate, endDate time.Time) (int, []time.Time) {
	if endDate.Before(startDate) {
		return 0, nil
	}

	current := normalizeDate(startDate)
	end := normalizeDate(endDate)
	workingDates := make([]time.Time, 0)

	for !current.After(end) {
		weekday := current.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			workingDates = append(workingDates, current)
		}
		current = current.AddDate(0, 0, 1)
	}
	return len(workingDates), workingDates
}

func CalculateAvailableBalance(total, reserved, pending, approved int) int {
	available := total - reserved - (pending + approved)
	if available < 0 {
		return 0
	}
	return available
}

func normalizeDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}
