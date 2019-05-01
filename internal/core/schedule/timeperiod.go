package schedule

// TimePeriod identifies which time period a schedule interval applies
type TimePeriod uint8

// TimePeriod constants define which time period a schedule interval applies
const (
	TimePeriodNone TimePeriod = iota
	TimePeriodHour
	TimePeriodDay
	TimePeriodWeek
	TimePeriodMonth
)

func (tp TimePeriod) String() string {
	switch tp {
	case TimePeriodNone:
		return "None"
	case TimePeriodHour:
		return "Hour"
	case TimePeriodDay:
		return "Day"
	case TimePeriodWeek:
		return "Week"
	case TimePeriodMonth:
		return "Month"
	}
	return "[Invalid time period]"
}
