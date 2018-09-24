package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("received request")
		name, _ := os.Hostname()
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		fmt.Fprintf(w, "Time: %s\nHost: %s\n  Ip: %s\n", time.Now(), name, ip)
	})
	log.Println("start server")
	server := &http.Server{Addr: ":8080"}
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
