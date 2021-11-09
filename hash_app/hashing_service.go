package hash_app

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

// regex needs to compiled at package level
var (
	re, _ = regexp.Compile("/hash/(.*)")
)

const (
	delay = 5 * time.Second
)

type hashingService interface {
	hash(id int64, hash string, startTime time.Time) int64
	hashWithDelay(hash string, startTime time.Time) int64
	getHash(id int64) string
	getStats() Stats
	getIdFromPath(path string) (int, error)
	reset()
}

type simpleHashingService struct {
	hashStore       hashStore
	aggregator      aggregator
	hashingFunction hashingFunction
}

func (h *simpleHashingService) hash(id int64, hash string, startTime time.Time) int64 {
	s := h.hashingFunction.hash(hash)
	h.hashStore.add(id, s)
	interval := time.Now().Sub(startTime)
	var ms = int64(time.Millisecond)
	h.aggregator.add(interval.Nanoseconds() / ms)
	return id
}

func (h *simpleHashingService) hashWithDelay(hash string, startTime time.Time) int64 {
	id := h.hashStore.getNextId()
	go func() {
		log.Printf("Hashing with delay %s \n", delay.String())
		time.Sleep(delay)
		h.hash(id, hash, startTime)
	}()
	return id
}

func (h *simpleHashingService) getHash(id int64) string {
	return h.hashStore.get(id)
}

func (h *simpleHashingService) getStats() Stats {
	return h.aggregator.get()
}

func (h *simpleHashingService) getIdFromPath(path string) (int, error) {
	match := re.FindStringSubmatch(path)
	if len(match) < 2 {
		return 0, fmt.Errorf("id not found")
	}

	id, err := strconv.Atoi(match[1])

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *simpleHashingService) reset() {
	h.hashStore.reset()
}

/*func main() {
	fmt.Println("Hashing Service")
	s:= getHashingService()

	//s.hashStore.idCounter=0
	//s.hashStore.hashTable=map[int]string{}

	s.hash("test1", time.Now())

	fmt.Println(s.getHash(1))
	fmt.Println(s.getStats())


    s.hashWithDelay("test", time.Now())

	time.Sleep(time.Second*10)

	fmt.Println(s.getStats())
}*/

/*func getHashingService() *simpleHashingService {
	return &simpleHashingService{
		hashStore: &memoryStore{idCounter: 0, hashTable: map[int]string{}},
		aggregator: &defaultAggregator{total: 0, totalTime: 0},
		hashingFunction: &shaHashingFunction{},
	}
}*/