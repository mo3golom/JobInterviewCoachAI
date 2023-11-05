package structs

type (
	BidirectionalMap[L, R comparable] struct {
		leftToRight map[L]R
		rightToLeft map[R]L
	}
)

func NewBidirectionalMap[L, R comparable](pairs []Pair[L, R]) *BidirectionalMap[L, R] {
	m := &BidirectionalMap[L, R]{
		leftToRight: map[L]R{},
		rightToLeft: map[R]L{},
	}

	for _, pair := range pairs {
		m.Add(pair.Left, pair.Right)
	}

	return m
}

func (m *BidirectionalMap[L, R]) Add(left L, right R) bool {
	_, leftExists := m.leftToRight[left]
	_, rightExists := m.rightToLeft[right]

	if leftExists || rightExists {
		return false
	}

	m.leftToRight[left] = right
	m.rightToLeft[right] = left

	return true
}

func (m *BidirectionalMap[L, R]) GetLeft(right R) (left L, ok bool) {
	l, ok := m.rightToLeft[right]
	return l, ok
}

func (m *BidirectionalMap[L, R]) GetRight(left L) (right R, ok bool) {
	r, ok := m.leftToRight[left]
	return r, ok
}
