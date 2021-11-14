package hash_app

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// regex needs to compiled at package level
var (
	re, _ = regexp.Compile("/hash/(.*)")
)

const (
	delay = 5 * time.Second
)

/**
Interface for all operations for hashing operations.
The rest controller/hashServer will use this interface,
to execute different operations for the app. This hashing service
by default injects memoryStore, shaHashingFunction for hashStore
and hashingFunction. But this can be easily modified by passing
different implementation for hashStore and shaHashingFunction. Thus,
we have modularized each and every dependency for this hashing service.
*/
type hashingService interface {
	hash(id int32, hash string, startTime time.Time) int32
	hashWithDelay(hash string, startTime time.Time) int32
	getHash(id int32) string
	getStats() Stats
	getIdFromPath(path string) (int, error)
	reset()
	finishAllTasks()
}

type simpleHashingService struct {
	hashStore       hashStore
	aggregator      statsCalculator
	hashingFunction hashingFunction
	waitGroup       *sync.WaitGroup
}

func (h *simpleHashingService) hash(id int32, hash string, startTime time.Time) int32 {
	s := h.hashingFunction.hash(hash)
	h.hashStore.add(id, s)
	interval := time.Now().Sub(startTime)
	var ms = int64(time.Millisecond)
	h.aggregator.add(int32(interval.Nanoseconds() / ms))
	return id
}

func (h *simpleHashingService) hashWithDelay(hash string, startTime time.Time) int32 {
	id := h.hashStore.getNextId()
	h.waitGroup.Add(1)
	go func() {
		defer h.waitGroup.Done()
		log.Printf("Will hash after delay of %s \n", delay.String())
		time.Sleep(delay)
		h.hash(id, hash, startTime)
	}()
	return id
}

func (h *simpleHashingService) getHash(id int32) string {
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

func (h *simpleHashingService) finishAllTasks() {
	h.waitGroup.Wait()
}