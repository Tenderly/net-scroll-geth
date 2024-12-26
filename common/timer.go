package common

import "github.com/tenderly/net-scroll-geth/metrics"

// WithTimer calculates the interval of f
func WithTimer(timer metrics.Timer, f func()) {
	if metrics.Enabled {
		timer.Time(f)
	} else {
		f()
	}
}
