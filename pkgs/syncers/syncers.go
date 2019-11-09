package syncers

type Syncer interface {
	Sync()
}

type syncer struct {
}

var _ Syncer = (*syncer)(nil)

func (s *syncer) Sync() {

}
