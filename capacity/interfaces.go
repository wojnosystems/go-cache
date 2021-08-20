package capacity

// Tracker tracks length and maxLen capacity
type Tracker interface {
	// IsLargerThanCapacity is true if the item will never fit into allotted capacity
	// false if it could fit as-is or if other items were removed
	IsLargerThanCapacity(itemSize uint) bool
}

type Mutator interface {
	// Add the amount to the internal state, false if cannot
	Add(amount uint) bool

	// Remove the amount, will not go below zero
	Remove(amount uint)
}

type TrackMutator interface {
	Tracker
	Mutator
}

// lenGetter obtains the current length of whatever is being tracked
type lenGetter func() uint
