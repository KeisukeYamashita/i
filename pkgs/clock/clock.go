package clock

import "time"

// Clock ...
type Clock interface {
	Now() time.Time
}

type clock struct {
	CurrentTime time.Time
}

var _ Clock = (*clock)(nil)

func NewClock(currentTime time.Time) Clock {
	return &clock{
		CurrentTime: currentTime,
	}
}

// Now ...
func (c *clock) Now() time.Time {
	return c.CurrentTime
}
