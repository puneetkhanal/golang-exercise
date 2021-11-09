package hash_app

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGetIdFromPath(t *testing.T) {
	hs := getHashingService()
	hs.reset()
	id, err := hs.getIdFromPath("/hash/1")

	assertEquals(t, id, 1)
	assertEquals(t, err, nil)

	id1, err := hs.getIdFromPath("/hash/jpt")

	assertEquals(t, id1, 0)
	assertEquals(t, err != nil, true)
}

func TestAddHashWithoutDelay(t *testing.T) {
	hs := getHashingService()
	hs.reset()
	hs.hash(1, "test", time.Now())

	id := hs.getHash(1)

	assertEquals(t, id, "7iaw3Ur350mqGo7jwQrpkj9hiYB3Lkc/iBml1JQODbJ6wYX4oOHV+E+IvIh/1nsUNzLDBMxfqa2Ob1f1ACio/w==")
}

func TestAddHashWithDelay(t *testing.T) {
	hs := getHashingService()
	hs.reset()
	hs.hashWithDelay("test", time.Now())

	// querying immediately should return empty
	id := hs.getHash(1)
	assertEquals(t, id, "")

	time.Sleep(10 * time.Second)

	id = hs.getHash(1)

	assertEquals(t, id, "7iaw3Ur350mqGo7jwQrpkj9hiYB3Lkc/iBml1JQODbJ6wYX4oOHV+E+IvIh/1nsUNzLDBMxfqa2Ob1f1ACio/w==")
}

func TestAggregator(t *testing.T) {
	hs := getHashingService()
	hs.reset()

	hs.hashWithDelay("test", time.Now())

	time.Sleep(5 * time.Second)

	hs.hashWithDelay("test1", time.Now())

	time.Sleep(15 * time.Second)

	fmt.Println(hs.getStats())

	stats := hs.getStats()
	assertEquals(t, stats.Total, int64(2))
	assertEquals(t, stats.Average >= 5000, true)
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
