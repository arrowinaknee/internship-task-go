package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type status int

const (
	status_ok          = status(0)
	status_waiting     = status(1)
	status_unavailable = status(2)
)

type sensor struct {
	host   string
	value  int
	status status
}

// should pass them procedurally
var sensors = []sensor{
	local(8081), local(8082), local(8083), local(8083),
}

// update sensor information every 30 seconds (plus request time)
func (s *sensor) observe() {
	for {
		s.status = status_waiting
		res, err := http.Get(s.host)
		if err != nil {
			log.Printf("error reaching host '%s': %s", s.host, err)
			s.status = status_unavailable
		} else {
			var data int
			_, err = fmt.Fscan(res.Body, &data)
			if err != nil {
				s.status = status_unavailable
			} else {
				s.value = data
				s.status = status_ok
			}
		}
		time.Sleep(30 * time.Second)
	}
}

// make sensor struct from port number on localhost
func local(port int) sensor {
	return sensor{host: fmt.Sprintf("http://localhost:%d", port)}
}

// convert this to return json
func handleHttp(w http.ResponseWriter, r *http.Request) {
	var resp struct {
		Value      int  `json:"value"`       // average sensor value
		IsOutdated bool `json:"is_outdated"` // some sensor data was not available
	}

	for _, s := range sensors {
		resp.Value += s.value
		resp.IsOutdated = resp.IsOutdated || s.status != status_ok
	}
	resp.Value /= len(sensors)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Panicf("error encoding response: %s", err)
	}
}

func main() {
	for i := range sensors {
		go sensors[i].observe()
	}

	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(handleHttp)))
}
