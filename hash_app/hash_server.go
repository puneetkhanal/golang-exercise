package hash_app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type hashServer struct {
	hashingService hashingService
	Server         *http.Server
	Shutdown       chan os.Signal
}

func (hs *hashServer) postHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				return
			}
			password, result := r.PostForm["password"]

			if !result {
				http.Error(w, "Password is required", http.StatusBadRequest)
				return
			}

			id := hs.hashingService.hashWithDelay(password[0], time.Now())
			fmt.Fprintf(w, "%d", id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (hs *hashServer) getHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id, err := hs.hashingService.getIdFromPath(r.URL.Path)

			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			hash := hs.hashingService.getHash(int32(id))
			fmt.Fprintf(w, "%s", hash)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (hs *hashServer) getStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(hs.hashingService.getStats())
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (hs *hashServer) getShutdownHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			hs.Stop()
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (hs *hashServer) start() error {
	log.Println("starting Server")
	return hs.Server.ListenAndServe()
}

// Stop ref: https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
func (hs *hashServer) Stop() {
	go func() {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		log.Println("Stopping Server and stop accepting new connections")

		if err := hs.Server.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
		log.Println("waiting for running task to complete")
		hs.hashingService.finishAllTasks()
		close(hs.Shutdown)
	}()
}

func GetHashServer(address string, port string) *hashServer {
	newHashServer := &hashServer{
		hashingService: getHashingService(),
		Shutdown: make(chan os.Signal, 1),
	}

	signal.Notify(newHashServer.Shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.Handle("/hash", newHashServer.postHandler())
	mux.Handle("/hash/", newHashServer.getHandler())
	mux.Handle("/stats", newHashServer.getStatsHandler())
	mux.Handle("/shutdown", newHashServer.getShutdownHandler())

	newHashServer.Server = &http.Server{Addr: address + ":" + port, Handler: mux}

	return newHashServer
}

func getHashingService() hashingService {
	return &simpleHashingService{
		hashStore:       &memoryStore{idCounter: 0, hashTable: map[int32]string{}},
		statsCalculator: &averageCalculator{total: 0, totalTime: 0},
		hashingFunction: &shaHashingFunction{},
		waitGroup:       &sync.WaitGroup{},
	}
}
