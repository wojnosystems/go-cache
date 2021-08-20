package capacity

type maxLen struct {
	capacity uint
	lenGet   lenGetter
}

func NewMaxLen(capacity uint, lenGet lenGetter) Tracker {
	return &maxLen{
		capacity: capacity,
		lenGet:   lenGet,
	}
}

func (m *maxLen) HasCapacity(itemSize uint) bool {
	return m.lenGet()+itemSize <= m.capacity
}

func (m *maxLen) IsLargerThanCapacity(itemSize uint) bool {
	return m.capacity < itemSize
}
