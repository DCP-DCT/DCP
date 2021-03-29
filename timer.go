package DCP

import "time"

type Timer struct {
	Timers map[string]time.Duration
}

func (t *Timer) Time(timerName string, from time.Time) {
	endTime := time.Now()
	t.Timers[timerName] = endTime.Sub(from)
}

func NewTimer(timerName string) (string, time.Time) {
	return timerName, time.Now()
}
