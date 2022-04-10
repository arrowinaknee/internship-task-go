package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// should pass them procedurally
var sensors = []string{
	"http://localhost:8081",
	"http://localhost:8082",
	"http://localhost:8083",
	"http://localhost:8084",
}

var sensor_vals []int    // should also srore status (pending, unavailable, etc.)
var data_lock sync.Mutex // is it needed here?

func observe(host string, id int) {
	for {
		res, err := http.Get("http://localhost:8081")
		if err != nil {
			log.Printf("error reaching host #%d: %s", id, err)
		} else {
			var data int
			fmt.Fscan(res.Body, &data)

			data_lock.Lock()
			sensor_vals[id] = data
			data_lock.Unlock()
		}
		time.Sleep(30 * time.Second)
	}
}

// convert this to return json
func handleHttp(w http.ResponseWriter, r *http.Request) {
	var sum = 0
	for _, v := range sensor_vals {
		sum += v
	}
	res := sum / len(sensors)
	fmt.Fprint(w, res)
}

func main() {
	sensor_vals = make([]int, len(sensors))

	for i, h := range sensors {
		go observe(h, i)
	}

	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(handleHttp)))
}
