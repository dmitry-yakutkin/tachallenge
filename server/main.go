package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tachallenge/server/set"
)

const (
	serverPort         = 8080
	httpRequestTimeout = 20

	linkKey = "u"
)

var (
	portStr string

	httpClient = &http.Client{Timeout: httpRequestTimeout * time.Second}
)

type NumbersServiceResponse struct {
	Numbers []int `json:"numbers"`
}

type ExternalNumbersServiceResponse struct {
	NumbersServiceResponse
}

func init() {
	portStr = fmt.Sprintf(":%v", serverPort)
}

func numbersHandler(w http.ResponseWriter, r *http.Request) {
	links := r.URL.Query()[linkKey]
	data := set.NewIntSet()

	for _, link := range links {
		func() {
			resp, err := httpClient.Get(link)
			if err != nil {
				return
			}
			defer r.Body.Close()

			respJSON := &ExternalNumbersServiceResponse{}
			err = json.NewDecoder(resp.Body).Decode(respJSON)
			if err != nil {
				return
			}
			data.Update(respJSON.Numbers)

		}()
	}
	enc := json.NewEncoder(w)
	enc.Encode(NumbersServiceResponse{data.Elements()})
}

func main() {
	http.HandleFunc("/numbers", numbersHandler)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
