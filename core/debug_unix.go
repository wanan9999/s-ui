//go:build linux

package core

import "syscall"

func rusageMaxRSS() float64 {
	var ru syscall.Rusage
	if syscall.Getrusage(syscall.RUSAGE_SELF, &ru) != nil {
		return 0
	}
	return float64(ru.Maxrss) / (1 << 10)
}
