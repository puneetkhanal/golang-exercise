package hash_app

import (
	"crypto/sha512"
	"encoding/base64"
)

type hashingFunction interface {
	hash(value string) string
}

type simpleHashingFunction struct {
}

type shaHashingFunction struct {
}

func (h *simpleHashingFunction) hash(value string) string {
	return value
}

func (h *shaHashingFunction) hash(value string) string {
	sha512Sum := sha512.Sum512([]byte(value))
	return base64.StdEncoding.EncodeToString(sha512Sum[:])
}

/*func main() {
	fmt.Println("Hashing Function - tests")
	s := simpleHashingFunction{}
	h := shaHashingFunction{}

	fmt.Println(s.hash("test"))
	fmt.Println(h.hash("test"))
}*/