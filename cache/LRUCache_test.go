package cache

import (
	"testing"
)

type Put struct {
	key   string
	value string
}

func (p Put) Appply(cache *LRUCache) {
	cache.Put(p.key, p.value)
}

type Get struct {
	key string
}

func (p Get) Appply(cache *LRUCache) {
	cache.Get(p.key)
}

func TestPutOneKey(t *testing.T) {

	type Operation interface {
		Appply(cache *LRUCache)
	}
	type node struct {
		key   string
		value string
	}

	tt := []struct {
		name       string
		size       int
		operations []Operation
		checkValues    map[string]string
		checkNoPresent []string
	}{
		{
			name:        "InsertOneValu",
			size:        3,
			operations:  []Operation{Put{"1", "a"}},
			checkValues: map[string]string{"1": "a"},
		},
		{
			name:           "Equal max size",
			size:           3,
			operations:     []Operation{Put{key: "1", value: "a"}, Put{key: "2", value: "a"}, Put{key: "3", value: "a"}},
			checkValues:    map[string]string{"1": "a", "2": "a", "3": "a"},
			checkNoPresent: []string{"4"},
		},
		{
			name:           "Surpas max size",
			size:           3,
			operations:     []Operation{Put{key: "1", value: "a"}, Put{key: "2", value: "a"}, Put{key: "3", value: "a"}, Put{key: "4", value: "a"}},
			checkValues:    map[string]string{"2": "a", "3": "a", "4": "a"},
			checkNoPresent: []string{"1"},
		},
		{
			name:        "Repeat keys",
			size:        3,
			operations:  []Operation{Put{key: "1", value: "a"}, Put{key: "2", value: "a"}, Put{key: "3", value: "a"}, Put{key: "1", value: "b"}},
			checkValues: map[string]string{"1": "b", "2": "a", "3": "a"},
		},
		{
			name:           "4 Put and 1 get",
			size:           3,
			operations:     []Operation{Put{"1", "a"}, Put{"2", "a"}, Put{"3", "a"}, Get{"1"}, Put{"4", "a"}},
			checkValues:    map[string]string{"1": "a", "3": "a", "4": "a"},
			checkNoPresent: []string{"2"},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			cache := NewLRUCache(test.size)
			for _, op := range test.operations {
				op.Appply(&cache)
			}

			for key, expectedValue := range test.checkValues {
				if value, ok := cache.Get(key); !ok {
					t.Errorf("Key %s not present in cache", key)
				} else if value != expectedValue {
					t.Errorf("Expected value %s is not equal to %s for key %s", expectedValue, value, key)
				}
			}

			for _, key := range test.checkNoPresent {
				if value, ok := cache.Get(key); ok {
					t.Errorf("Key %s was found in cache with value %s", key, value)
				}
			}
		})
	}
}
