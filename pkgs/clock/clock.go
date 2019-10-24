package clock

import "time"

// Clock ...
type Clock interface {
	Now() time.Time
}

type clock struct {
	CreatedAt time.Time
}

var _ Clock = (*clock)(nil)

// NewClock ...
func NewClock(currentTime time.Time) Clock {
	return &clock{
		CreatedAt: currentTime,
	}
}

// Now ...
func (c *clock) Now() time.Time {
	return time.Now()
}
