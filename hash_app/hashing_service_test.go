package hash_app

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGetIdFromPath(t *testing.T) {
	hashingService := getHashingService()
	hashingService.reset()
	id, err := hashingService.getIdFromPath("/hash/1")

	assertEquals(t, id, 1)
	assertEquals(t, err, nil)

	id, err = hashingService.getIdFromPath("/hash/jpt")

	assertEquals(t, id, 0)
	assertEquals(t, err != nil, true)

	id, err = hashingService.getIdFromPath("/hash")
	assertEquals(t, id, 0)
	assertEquals(t, err.Error(), "id not found")
}

func TestAddHashWithoutDelay(t *testing.T) {
	hashingService := getHashingService()
	hashingService.reset()
	hashingService.hash(1, "test", time.Now())

	id := hashingService.getHash(1)

	assertEquals(t, id, "7iaw3Ur350mqGo7jwQrpkj9hiYB3Lkc/iBml1JQODbJ6wYX4oOHV+E+IvIh/1nsUNzLDBMxfqa2Ob1f1ACio/w==")
}

func TestAddHashWithDelay(t *testing.T) {
	hashingService := getHashingService()
	hashingService.reset()
	hashingService.hashWithDelay("test", time.Now())

	// querying immediately should return empty
	id := hashingService.getHash(1)
	assertEquals(t, id, "")

	time.Sleep(10 * time.Second)

	id = hashingService.getHash(1)

	assertEquals(t, id, "7iaw3Ur350mqGo7jwQrpkj9hiYB3Lkc/iBml1JQODbJ6wYX4oOHV+E+IvIh/1nsUNzLDBMxfqa2Ob1f1ACio/w==")
}

func TestAggregator(t *testing.T) {
	hashingService := getHashingService()
	hashingService.reset()

	hashingService.hashWithDelay("test", time.Now())

	time.Sleep(5 * time.Second)

	hashingService.hashWithDelay("test1", time.Now())

	time.Sleep(15 * time.Second)

	fmt.Println(hashingService.getStats())

	stats := hashingService.getStats()
	assertEquals(t, stats.Total, int32(2))
	assertEquals(t, stats.Average >= float32(5000), true)
}

// This seems to best approach for fluent assertions.
// ref: https://gist.github.com/samalba/6059502
// AssertEqual checks if values are equal
func assertEquals(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
