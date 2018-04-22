package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dmitry-yakutkin/tachallenge/server/fetch"
	"github.com/dmitry-yakutkin/tachallenge/server/set"
)

type NumbersServiceResponse struct {
	Numbers []int `json:"numbers"`
}

func processLinks(fetcher fetch.Fetcher, links []string, numbers *[]int) {
	timer := time.NewTimer(maxRequestProcessingDuration)
	data := set.NewIntSet()
	doneFetchingNumbers := make(chan bool)

	// Links processing goroutines.
	for _, link := range links {
		go func(link string) {
			resp, err := fetcher.Get(link)
			if err != nil {
				log.Printf("%s fetching has failed", link)
				doneFetchingNumbers <- false
				return
			}
			defer resp.Body.Close()

			respJSON := &NumbersServiceResponse{}
			err = json.NewDecoder(resp.Body).Decode(respJSON)
			if err != nil {
				log.Printf("%s decoding has failed", link)
				doneFetchingNumbers <- false
				return
			}

			data.Update(respJSON.Numbers)
			log.Printf("%s is processed successfully", link)
			doneFetchingNumbers <- true
		}(link)
	}

	processedLinks := 0

	// Links processing goroutines management: exit either becasuse of all links where processed or
	// timeout was exceeded.
OUT:
	for {
		select {
		case <-doneFetchingNumbers:
			processedLinks++
			if processedLinks == len(links) {
				log.Printf("all links are processed successfully")
				break OUT
			}
		case <-timer.C:
			log.Printf("processing timeout, stopped processing")
			break OUT
		}
	}

	// Return resulting numbers back.
	elements := data.Elements()
	*numbers = append(*numbers, elements...)
}

func numbersHandler(w http.ResponseWriter, r *http.Request) {
	links := r.URL.Query()[linkKey]
	log.Printf("processing %v.", links)

	var numbers []int
	processLinks(fetch.NewHTTPFetcher(), links, &numbers)

	enc := json.NewEncoder(w)
	log.Printf("finished processing, results: %v.", numbers)
	enc.Encode(NumbersServiceResponse{numbers})
}
