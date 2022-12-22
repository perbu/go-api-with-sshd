package api

func (c *Counters) IncErrors() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errors++
}

func (c *Counters) IncRequests() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.requests++
}

func (c *Counters) GetErrors() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.errors
}

func (c *Counters) GetRequests() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.requests
}
