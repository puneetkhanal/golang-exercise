package hash_app

type aggregator interface {
	add(totalTime int64)
	get() Stats
}

type defaultAggregator struct {
	total     int64
	totalTime int64
}

type Stats struct {
	Total   int64 `json:"total"`
	Average int64 `json:"average"`
}

func (r *defaultAggregator) add(totalTime int64) {
	r.total = r.total + 1
	r.totalTime = r.totalTime + totalTime
}

func (r *defaultAggregator) get() Stats {
	if r.total == 0 {
		return Stats{Total: 0, Average: 0}
	} else {
		return Stats{Total: r.total, Average: r.totalTime / r.total}
	}
}

/*func main() {
	fmt.Println("Aggregator Test")
	r := defaultAggregator{ total: 0, totalTime: 0}
	r.add(1)
	fmt.Println(r.total)
	fmt.Println(r.totalTime)
}*/
