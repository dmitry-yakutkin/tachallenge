package set

import (
	"sort"
	"sync"
)

type IntSet struct {
	SetMap *sync.Map
}

func NewIntSet() IntSet {
	return IntSet{SetMap: &sync.Map{}}
}

func (s IntSet) Set(value int) {
	s.SetMap.Store(value, true)
}

func (s IntSet) Update(values []int) {
	for _, value := range values {
		s.Set(value)
	}
}

func (s IntSet) Elements() []int {
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
