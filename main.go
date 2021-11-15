package main

import (
	hashApp "hash/app/hash_app"
	"log"
	"syscall"
)

func main() {
	hashServer := hashApp.GetHashServer("localhost", "9090")

	// Run hashServer in the background
	go func() {
		if err := hashServer.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Print("Server Started")

	// block main thread until channel is closed by invoking /shutdown api or syscall.SIGINT and syscall.SIGTERM
	ch := <-hashServer.Shutdown
	if ch == syscall.SIGINT || ch == syscall.SIGTERM {
		hashServer.Stop()
		<-hashServer.Shutdown
	}
	log.Print("Server Exited Properly")
}
