package pokecache

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val LocationResponse
	}{
		{
			key: "https://example.com",
			val: LocationResponse{
				Results: []Location{
					{
						Name: "testdata",
					},
				},
			},
		},
		{
			key: "https://example.com/path",
			val: LocationResponse{
				Results: []Location{
					{
						Name: "moretestdata",
					},
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)

			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key %v in cache", c.key)
				return
			}

			// Use reflect.DeepEqual to compare slices
			if !reflect.DeepEqual(val, c.val) {
				t.Errorf("expected value %v, but got %v", c.val, val)
				return
			}
		})
	}
}

func TestDelete(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", LocationResponse{
		Results: []Location{
			{
				Name: "testdata",
			},
		},
	})

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
