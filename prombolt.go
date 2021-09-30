// Package prombolt provides a Prometheus metrics collector for Bolt databases.
package prombolt

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	bolt "go.etcd.io/bbolt"
)

const (
	// namespace is the top-level namespace metric names.
	namespace = "bolt"
)

// New creates a new prometheus.Collector that can be registered with
// Prometheus to scrape metrics from a Bolt database handle.
//
// Name should specify a unique name for the collector, and will be added
// as a label to all produced Prometheus metrics.
func New(name string, db *bolt.DB) prometheus.Collector {
	return &collector{
		stats:       newStatsCollector(name, db),
		bucketStats: newBucketStatsCollector(name, db),
	}
}

// Enforce that collector is a prometheus.Collector.
var _ prometheus.Collector = &collector{}

// A collector is a prometheus.Collector for Bolt database metrics.
type collector struct {
	mu          sync.Mutex
	stats       *statsCollector
	bucketStats *bucketStatsCollector
}

// Describe implements the prometheus.Collector interface.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats.Describe(ch)
	c.bucketStats.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats.Collect(ch)
	c.bucketStats.Collect(ch)
}
