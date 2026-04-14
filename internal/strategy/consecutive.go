package strategy

type Consecutive struct {
	failures  int
	threshold int
}

func NewConsecutive(threshold int) *Consecutive {
	return &Consecutive{threshold: threshold}
}

func (c *Consecutive) Record(success bool) {
	if success {
		c.failures = 0
	} else {
		c.failures++
	}
}

func (c *Consecutive) ShouldOpen() bool {
	return c.failures >= c.threshold
}

func (c *Consecutive) Reset() {
	c.failures = 0
}
