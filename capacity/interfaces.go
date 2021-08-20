package capacity

// Tracker tracks length and maxLen capacity
type Tracker interface {
	// HasCapacity true if it can fit with existing items, false if not
	HasCapacity(itemSize uint) bool

	// IsLargerThanCapacity is true if the item will never fit into allotted capacity
	// false if it could fit as-is or if other items were removed
	IsLargerThanCapacity(itemSize uint) bool
}

// lenGetter obtains the current length of whatever is being tracked
type lenGetter func() uint
