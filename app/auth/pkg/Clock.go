package pkg

import "time"

type Clock interface {
	Now() time.Time
}

type NormalClock struct {
}

func (n NormalClock) Now() time.Time {
	return time.Now()
}

type StubClock struct {
	time.Time
}

func (s StubClock) Now() time.Time {
	return s.Time
}
