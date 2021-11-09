package hash_app

type hashStore interface {
	add(id int64, hash string) int64
	get(id int64) string
	getNextId() int64
	reset()
}

type memoryStore struct {
	idCounter int64            `default:"1"`
	hashTable map[int64]string `default:"{}"`
}

func (h *memoryStore) add(id int64, hash string) int64 {
	h.hashTable[id] = hash
	return id
}

func (h *memoryStore) get(id int64) string {
	return h.hashTable[id]
}

func (h *memoryStore) getNextId() int64 {
	h.idCounter++
	return h.idCounter
}

func (h *memoryStore) reset() {
	h.idCounter = 0
}

/*func main() {
	fmt.Println("Hash Store tests")
	r := memoryStore{ idCounter: 0, hashTable: map[int]string{}}
	r.add("test")
	r.add("test1")
	fmt.Println(r.hashTable)
	fmt.Println(r.get(1))
	fmt.Println(r.get(2))
}*/
