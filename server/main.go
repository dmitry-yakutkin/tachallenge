package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	serverPort         = 8080
	httpRequestTimeout = 20

	maxRequestProcessingDuration       = time.Millisecond * 500
	maxExternalNumbersFetchingDuration = time.Millisecond * 400

	linkKey = "u"
)

var (
	portStr string
)

func init() {
	portStr = fmt.Sprintf(":%v", serverPort)
}

func main() {
	http.HandleFunc("/numbers", numbersHandler)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
