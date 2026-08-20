package clock

import "time"

// Fake 固定/可推进时钟。
type Fake struct {
	T time.Time
}

func (f *Fake) Now() time.Time {
	if f.T.IsZero() {
		return time.Unix(0, 0).UTC()
	}
	return f.T.UTC()
}

// Advance 推进。
func (f *Fake) Advance(d time.Duration) {
	if f.T.IsZero() {
		f.T = time.Unix(0, 0).UTC()
	}
	f.T = f.T.Add(d)
}
