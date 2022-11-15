package api

func (c *counters) IncErrors() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errors++
}

func (c *counters) IncRequests() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.requests++
}

func (c *counters) GetErrors() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.errors
}

func (c *counters) GetRequests() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.requests
}
