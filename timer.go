package DCP

import "time"

type TimerEntry struct {
	Avg     int64 `json:"avg"`
	Counter int   `json:"counter"`
}

type Timer struct {
	Timers map[string]TimerEntry `json:"timers"`
}

func (t *Timer) Time(timerName string, from time.Time) {
	endTime := time.Now()

	timerEntry := TimerEntry{
		Avg:     0,
		Counter: 1,
	}

	if te, exist := t.Timers[timerName]; exist {
		timerEntry = te
	}

	runTime := endTime.Sub(from)

	timerEntry.Avg = (timerEntry.Avg + int64(runTime)) / int64(timerEntry.Counter)
	timerEntry.Counter = timerEntry.Counter + 1

	t.Timers[timerName] = timerEntry
}

func NewTimer(timerName string) (string, time.Time) {
	return timerName, time.Now()
}
