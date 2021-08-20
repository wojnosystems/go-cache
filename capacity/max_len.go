package capacity

type maxLen struct {
	cap uint
	len uint
}

func NewMaxLen(cap uint) TrackMutator {
	return &maxLen{
		cap: cap,
	}
}

func (m *maxLen) IsLargerThanCapacity(itemSize uint) bool {
	return m.cap < itemSize
}

func (m *maxLen) Add(amount uint) bool {
	if !m.hasCapacity(amount) {
		return false
	}
	m.len += amount
	return true
}

func (m *maxLen) hasCapacity(itemSize uint) bool {
	return m.len+itemSize <= m.cap
}

func (m *maxLen) Remove(amount uint) {
	if m.len < amount {
		m.len = 0
	} else {
		m.len = m.len - amount
	}
}
