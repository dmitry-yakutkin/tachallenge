package set

import (
	"sort"
	"sync"
)

// IntSet is a basic implementation of Set data sctructure, which operates with int values.
type IntSet interface {
	// Set adds value to the set.
	Set(value int)

	// Update adds multiple values to the set.
	Update(values []int)

	// Elements retreives all elements currently added to the set.
	Elements() []int
}

type intSet struct {
	SetMap *sync.Map
}

// NewIntSet constructs IntSet instances.
func NewIntSet() IntSet {
	return intSet{SetMap: &sync.Map{}}
}

func (s intSet) Set(value int) {
	s.SetMap.Store(value, true)
}

func (s intSet) Update(values []int) {
	for _, value := range values {
		s.Set(value)
	}
}

func (s intSet) Elements() []int {
	result := []int{}
	s.SetMap.Range(func(key, value interface{}) bool {
		intKey, ok := key.(int)
		if !ok {
			return false
		}
		result = append(result, intKey)
		return true
	})
	sort.Ints(result)
	return result
}
