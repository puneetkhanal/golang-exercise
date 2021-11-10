package hash_app

/**
Interface for statsCalculator.
Provides implementation for default averageCalculator.
Other different statsCalculator can be implemented later.
*/
type statsCalculator interface {
	add(totalTime int64)
	get() Stats
}

type averageCalculator struct {
	total     int64
	totalTime int64
}

type Stats struct {
	Total   int64 `json:"total"`
	Average int64 `json:"average"`
}

func (r *averageCalculator) add(totalTime int64) {
	r.total = r.total + 1
	r.totalTime = r.totalTime + totalTime
}

func (r *averageCalculator) get() Stats {
	if r.total == 0 {
		return Stats{Total: 0, Average: 0}
	} else {
		return Stats{Total: r.total, Average: r.totalTime / r.total}
	}
}

/*func main() {
	fmt.Println("Aggregator Test")
	r := averageCalculator{ total: 0, totalTime: 0}
	r.add(1)
	fmt.Println(r.total)
	fmt.Println(r.totalTime)
}*/
