package lru

// Tracker of least recently used key
type Tracker interface {
	// Touch the key, create or mark existing as recently used
	Touch(key interface{})

	// Remove the key provided
	Remove(key interface{})

	// LRU get the least recently used key
	LRU() (key interface{}, ok bool)

	// Len how many items tracked in this structure
	Len() int
}
