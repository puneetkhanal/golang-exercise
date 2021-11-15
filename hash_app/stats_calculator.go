package hash_app

import "sync"

/**
Interface for statsCalculator.
Provides implementation for default averageCalculator.
Other different statsCalculator can be implemented later.
*/
type statsCalculator interface {
	add(totalTime int32)
	get() Stats
}

type averageCalculator struct {
	total     int32
	totalTime int32
	statsLock sync.RWMutex
}

type Stats struct {
	Total   int32   `json:"total"`
	Average float32 `json:"average"`
}

func (r *averageCalculator) add(totalTime int32) {
	r.statsLock.Lock()
	defer r.statsLock.Unlock()
	r.total += 1
	r.totalTime += totalTime
}

func (r *averageCalculator) get() Stats {
	r.statsLock.RLock()
	defer r.statsLock.RUnlock()
	if r.total == 0 {
		return Stats{Total: 0, Average: float32(0)}
	} else {
		return Stats{Total: r.total, Average: float32(r.totalTime) / float32(r.total)}
	}
}
