package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

type numbersEndpoint struct {
	initialData []int
	server      *httptest.Server
}

func genNumbersEndpoint(items []int) numbersEndpoint {
	return numbersEndpoint{
		initialData: items,
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			initialDataJSON := NumbersServiceResponse{items}
			initialDataBytes, _ := json.Marshal(initialDataJSON)
			fmt.Fprintln(w, string(initialDataBytes))
		})),
	}
}

func sortNums(numbers []int) []int {
	m := make(map[int]bool)
	for _, item := range numbers {
		m[item] = true
	}

	res := []int{}

	for key := range m {
		res = append(res, key)
	}

	sort.Ints(res)

	return res
}

func TestBasicProcessLinks(t *testing.T) {
	expectedNumbers := []int{}

	// Prepare testing endpoints.
	numbersEndpoints := []numbersEndpoint{
		genNumbersEndpoint([]int{1, 1, 2, 3, 5, 8, 13, 21}),
		genNumbersEndpoint([]int{5, 17, 3, 19, 76, 24, 1, 5, 10, 34, 8, 27, 7}),
	}

	// Make sure testing servers will be terminated at the end.
	defer func() {
		for _, c := range numbersEndpoints {
			c.server.Close()
		}
	}()

	queryStr := ""

	for _, c := range numbersEndpoints {
		queryStr += fmt.Sprintf("u=%s&", c.server.URL)
		expectedNumbers = append(expectedNumbers, c.initialData...)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", fmt.Sprintf("/numbers?%s", queryStr), nil)
	numbersHandler(w, r)

	resp := w.Result()

	actualDataJSON := NumbersServiceResponse{}
	json.NewDecoder(resp.Body).Decode(&actualDataJSON)

	expectedNumbers = sortNums(expectedNumbers)

	if !reflect.DeepEqual(expectedNumbers, actualDataJSON.Numbers) {
		t.Fatalf("expected %v, got %v", expectedNumbers, actualDataJSON.Numbers)
	}
}
