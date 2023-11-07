package main

import (
	"fmt"
	"github.com/prometheus-community/pro-bing"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func handlePing(w http.ResponseWriter, r *http.Request) {

	// HOST
	urlParam_host := r.URL.Query().Get("host")
	pinger, err := probing.NewPinger(urlParam_host)
	if err != nil {
		panic(err)
	}

	// TIMEOUT
	urlParam_timeout := r.URL.Query().Get("timeout")
	if urlParam_timeout == "" {
		urlParam_timeout = "1s"
	}
	pinger.Timeout, err = time.ParseDuration(urlParam_timeout)
	if err != nil {
		panic(err)
	}

	// COUNT
	urlParam_count := r.URL.Query().Get("count")
	if urlParam_count == "" {
		urlParam_count = "1"
	}
	count, err := strconv.Atoi(urlParam_count)
	if err != nil {
		panic(err)
	}
	if count > 0 {
		pinger.Count = count
	} else {
		pinger.Count = 1
	}

	// RUN PING
	err = pinger.Run()
	if err != nil {
		panic(err)
	}

	// OUTPUT
	if pinger.Statistics().PacketLoss == 100 {
		http.Error(w, "The requested host does not respond", http.StatusBadRequest)
	} else {
		s := fmt.Sprintf("%d  packets transmitted. %d received, %.2f%% packet loss\n",
			pinger.Statistics().PacketsSent,
			pinger.Statistics().PacketsRecv,
			pinger.Statistics().PacketLoss)
		w.Write([]byte(string(s)))
	}
}

func getListenPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return ":8080"
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlePing(w, r)
	})

	port := getListenPort()
	log.Printf("httpGoPing started listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
