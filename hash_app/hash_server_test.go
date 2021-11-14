package hash_app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGetHashByIdNotPresent(t *testing.T) {
	hashServer := GetHashServer("localhost", "8080")
	req := httptest.NewRequest("GET", "/hash/1", nil)
	w := httptest.NewRecorder()

	getHashHandler := hashServer.getHandler()
	getHashHandler(w, req)

	assertEquals(t, w.Code, http.StatusOK)
	assertEquals(t, w.Body.String(), "")
}

func TestCreateAndGetHash(t *testing.T) {
	hashServer := GetHashServer("localhost", "8080")
	postHashData := url.Values{}
	postHashData.Add("password", "angryMonkey")

	req := httptest.NewRequest("POST", "/hash", strings.NewReader(postHashData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	postHashHandler := hashServer.postHandler()
	postHashHandler(w, req)

	assertEquals(t, w.Code, http.StatusOK)
	assertEquals(t, w.Body.String(), "1")

	time.Sleep(7 * time.Second)
	req = httptest.NewRequest("GET", "/hash/1", nil)
	w = httptest.NewRecorder()

	getHashHandler := hashServer.getHandler()
	getHashHandler(w, req)

	assertEquals(t, w.Code, http.StatusOK)
	assertEquals(t, w.Body.String(), "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==")
}

func TestMissingPasswordParameter(t *testing.T) {
	hashServer := GetHashServer("localhost", "8080")
	postHashData := url.Values{}
	postHashData.Add("x", "y")

	req := httptest.NewRequest("POST", "/hash", strings.NewReader(postHashData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	postHashHandler := hashServer.postHandler()
	postHashHandler(w, req)

	assertEquals(t, w.Code, http.StatusBadRequest)
	assertEquals(t, w.Body.String(), "Password is required\n")
}

func TestGetStatsNoData(t *testing.T) {
	hashServer := GetHashServer("localhost", "8080")
	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()

	getStatsHandler := hashServer.getStatsHandler()
	getStatsHandler(w, req)

	data := Stats{}
	json.Unmarshal([]byte(w.Body.String()), &data)

	assertEquals(t, data.Average, float32(0))
	assertEquals(t, data.Total, int32(0))
}

func TestShutDown(t *testing.T) {
	hashServer := GetHashServer("localhost", "8080")

	req := httptest.NewRequest("GET", "/shutdown", nil)
	w := httptest.NewRecorder()

	getShutdownHandler := hashServer.getShutdownHandler()
	getShutdownHandler(w, req)

	assertEquals(t, w.Code, http.StatusOK)

	time.Sleep(5 * time.Second)

	select {
	case <-hashServer.Shutdown:
	default:
		t.Errorf("shutdown channel should be closed")
	}
}
