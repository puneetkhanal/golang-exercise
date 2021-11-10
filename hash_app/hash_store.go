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
	add(id int64, hash string) int64
	get(id int64) string
	getNextId() int64
	reset()
}

type memoryStore struct {
	idLock    sync.RWMutex
	writeLock sync.RWMutex
	idCounter int64            `default:"1"`
	hashTable map[int64]string `default:"{}"`
}

func (h *memoryStore) add(id int64, hash string) int64 {
	h.writeLock.Lock()
	defer h.writeLock.Unlock()
	//log.Printf("hash %s\n", hash)
	h.hashTable[id] = hash
	return id
}

func (h *memoryStore) get(id int64) string {
	return h.hashTable[id]
}

func (h *memoryStore) getNextId() int64 {
	h.idLock.Lock()
	defer h.idLock.Unlock()
	h.idCounter++
	return h.idCounter
}

func (h *memoryStore) reset() {
	h.idCounter = 0
}