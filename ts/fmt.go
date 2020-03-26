package ts

import (
	"fmt"
	"time"

	"github.com/willf/pad"
)

const (
	layoutS  = "2006-01-02 15:04:05"
	layoutMs = "2006-01-02 15:04:05.999"
)

func FromNano(ts int64) time.Time {
	return time.Unix(ts/int64(1e9), ts%int64(1e9))
}

func FromMilli(ts int64) time.Time {
	return time.Unix(ts/1000, (ts%1000)*int64(time.Millisecond))
}

func FmtTimeS(ts time.Time) string {
	return ts.Format(layoutS)
}

func FmtTimeMs(ts time.Time) string {
	return pad.Right(ts.Format(layoutMs), 23, "0")
}

func FmtDuration(v time.Duration) string {
	v = v.Round(time.Second)

	d := v / (24 * time.Hour)
	v -= d * 24 * time.Hour

	h := v / time.Hour
	v -= h * time.Hour

	m := v / time.Minute
	v -= m * time.Minute

	s := v / time.Second

	if d > 0 {
		return fmt.Sprintf("%02dd%02dh%02dm%02ds", d, h, m, s)
	}
	if h > 0 {
		return fmt.Sprintf("%02dh%02dm%02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%02dm%02ds", m, s)
	}

	return fmt.Sprintf("%02ds", s)
}
