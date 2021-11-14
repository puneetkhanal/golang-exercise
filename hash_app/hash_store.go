package hash_app

import (
	"sync"
)

/**
Interface for hashStore. Provides default implementation
for memoryStore. If we need an interface to other external db,
we can implement this interface e.g. RedisStore, PostgresStore, etc.
*/
type hashStore interface {
	add(id int32, hash string) int32
	get(id int32) string
	getNextId() int32
	reset()
}

type memoryStore struct {
	lock      sync.RWMutex
	idCounter int32            `default:"1"`
	hashTable map[int32]string `default:"{}"`
}

func (h *memoryStore) add(id int32, hash string) int32 {
	h.lock.Lock()
	defer h.lock.Unlock()
	//log.Printf("hash %s\n", hash)
	h.hashTable[id] = hash
	return id
}

func (h *memoryStore) get(id int32) string {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.hashTable[id]
}

func (h *memoryStore) getNextId() int32 {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.idCounter++
	return h.idCounter
}

func (h *memoryStore) reset() {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.idCounter = 0
}
