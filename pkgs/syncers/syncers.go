package syncers

import "time"

type Syncer interface {
	Sync()
}

type syncer struct {
	ticker *time.Ticker
}

var _ Syncer = (*syncer)(nil)

func (s *syncer) Sync() {
}
