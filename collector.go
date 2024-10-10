package libagent

type Collector struct {
	Key     string
	Version string
}

// Returns a new Collector
func NewCollector(key string) *Collector {
	return &Collector{
		Key: key,
	}
}
