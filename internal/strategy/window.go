package strategy

type SlidingWindow struct {
	size        int
	values      []bool
	index       int
	count       int
	failureRate float64
	minRequests int
}

func NewSlidingWindow(size int, rate float64, min int) *SlidingWindow {
	return &SlidingWindow{
		size:        size,
		values:      make([]bool, size),
		failureRate: rate,
		minRequests: min,
	}
}

func (w *SlidingWindow) Record(success bool) {
	w.values[w.index] = success
	w.index = (w.index + 1) % w.size

	if w.count < w.size {
		w.count++
	}
}

func (w *SlidingWindow) ShouldOpen() bool {
	if w.count < w.minRequests {
		return false
	}

	failures := 0
	for i := 0; i < w.count; i++ {
		if !w.values[i] {
			failures++
		}
	}

	rate := float64(failures) / float64(w.count)
	return rate >= w.failureRate
}

func (w *SlidingWindow) Reset() {
	w.count = 0
	w.index = 0
}
