package prombolt

import (
	"github.com/prometheus/client_golang/prometheus"
	bolt "go.etcd.io/bbolt"
)

var _ prometheus.Collector = &bucketStatsCollector{}

// A bucketStatsCollector is a prometheus.Collector for Bolt database bucket
// statistics.
type bucketStatsCollector struct {
	name    string
	db      *bolt.DB
	forEach func(fn forEachBucketStatsFunc) error

	LogicalBranchPages                *prometheus.Desc
	PhysicalBranchOverflowPages       *prometheus.Desc
	LogicalLeafPages                  *prometheus.Desc
	PhysicalLeafOverflowPages         *prometheus.Desc
	Keys                              *prometheus.Desc
	Depth                             *prometheus.Desc
	PhysicalBranchPagesAllocatedBytes *prometheus.Desc
	PhysicalBranchPagesInUseBytes     *prometheus.Desc
	PhysicalLeafPagesAllocatedBytes   *prometheus.Desc
	PhysicalLeafPagesInUseBytes       *prometheus.Desc
	Buckets                           *prometheus.Desc
	InlinedBuckets                    *prometheus.Desc
	InlinedBucketsInUseBytes          *prometheus.Desc
}

// newBucketStatsCollector creates a new bucketStatsCollector with the specified
// name and forEachBucketFunc for retrieving statistics.
func newBucketStatsCollector(name string, db *bolt.DB) *bucketStatsCollector {
	const (
		subsystem = "bucket"
	)

	var (
		labels = []string{"database", "bucket"}
	)

	return &bucketStatsCollector{
		name: name,
		db:   db,
		// By default, forEach iterates each bucket retrieved from the Bolt
		// database handle, but this is swappable for tests
		forEach: forEachWithBoltDB(db),

		LogicalBranchPages: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "logical_branch_pages"),
			"Number of logical branch pages for a bucket.",
			labels,
			nil,
		),

		PhysicalBranchOverflowPages: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "physical_branch_overflow_pages"),
			"Number of physical branch overflow pages for a bucket.",
			labels,
			nil,
		),

		LogicalLeafPages: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "logical_leaf_pages"),
			"Number of logical leaf pages for a bucket.",
			labels,
			nil,
		),

		PhysicalLeafOverflowPages: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "physical_leaf_overflow_pages"),
			"Number of physical leaf overflow pages for a bucket.",
			labels,
			nil,
		),

		Keys: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "keys"),
			"Number of key/value pairs in a bucket.",
			labels,
			nil,
		),

		Depth: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "depth"),
			"Number of levels in B+ tree for a bucket.",
			labels,
			nil,
		),

		PhysicalBranchPagesAllocatedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "physical_branch_pages_allocated_bytes"),
			"Number of bytes allocated in physical branch pages for a bucket.",
			labels,
			nil,
		),

		PhysicalBranchPagesInUseBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "physical_branch_pages_in_use_bytes"),
			"Number of bytes in use in physical branch pages for a bucket.",
			labels,
			nil,
		),

		PhysicalLeafPagesAllocatedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "physical_leaf_pages_allocated_bytes"),
			"Number of bytes allocated in physical leaf pages for a bucket.",
			labels,
			nil,
		),

		PhysicalLeafPagesInUseBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "physical_leaf_pages_in_use_bytes"),
			"Number of bytes in use in physical leaf pages for a bucket.",
			labels,
			nil,
		),

		Buckets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "buckets"),
			"Number of buckets within a bucket, including the top bucket.",
			labels,
			nil,
		),

		InlinedBuckets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "inlined_buckets"),
			"Number of inlined buckets for a bucket.",
			labels,
			nil,
		),

		InlinedBucketsInUseBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "inlined_buckets_in_use_bytes"),
			"Number of bytes in use for inlined buckets.",
			labels,
			nil,
		),
	}
}

// Describe implements the prometheus.Collector interface.
func (c *bucketStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.LogicalBranchPages,
		c.PhysicalBranchOverflowPages,
		c.LogicalLeafPages,
		c.PhysicalLeafOverflowPages,
		c.Keys,
		c.Depth,
		c.PhysicalBranchPagesAllocatedBytes,
		c.PhysicalBranchPagesInUseBytes,
		c.PhysicalLeafPagesAllocatedBytes,
		c.PhysicalLeafPagesInUseBytes,
		c.Buckets,
		c.InlinedBuckets,
		c.InlinedBucketsInUseBytes,
	}

	for _, d := range ds {
		ch <- d
	}
}

// A forEachBucketStatsFunc is a function which is repeatedly called for all
// buckets in a Bolt database to collect bucket statistics.
type forEachBucketStatsFunc func(bucket string, s bolt.BucketStats) error

// forEachWithBoltDB begins a read-only bolt transaction and returns a forEach
// function for a bucketStatsCollector.  The returned function is invoked
// repeatedly for each bucket and its stats retrieved from the Bolt database
// handle.
func forEachWithBoltDB(db *bolt.DB) func(forEachBucketStatsFunc) error {
	return func(iter forEachBucketStatsFunc) error {
		return db.View(func(tx *bolt.Tx) error {
			return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
				// TODO(mdlayher): if/when possible, iterate child buckets and
				// collect metrics for them as well.
				// See: https://github.com/boltdb/bolt/issues/603.
				return iter(string(name), b.Stats())
			})
		})
	}
}

// Collect implements the prometheus.Collector interface.
func (c *bucketStatsCollector) Collect(ch chan<- prometheus.Metric) {
	err := c.forEach(func(bucket string, s bolt.BucketStats) error {
		ch <- prometheus.MustNewConstMetric(
			c.LogicalBranchPages,
			prometheus.GaugeValue,
			float64(s.BranchPageN),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PhysicalBranchOverflowPages,
			prometheus.GaugeValue,
			float64(s.BranchOverflowN),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.LogicalLeafPages,
			prometheus.GaugeValue,
			float64(s.LeafPageN),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PhysicalLeafOverflowPages,
			prometheus.GaugeValue,
			float64(s.LeafOverflowN),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Keys,
			prometheus.GaugeValue,
			float64(s.KeyN),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Depth,
			prometheus.GaugeValue,
			float64(s.Depth),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PhysicalBranchPagesAllocatedBytes,
			prometheus.GaugeValue,
			float64(s.BranchAlloc),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PhysicalBranchPagesInUseBytes,
			prometheus.GaugeValue,
			float64(s.BranchInuse),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PhysicalLeafPagesAllocatedBytes,
			prometheus.GaugeValue,
			float64(s.LeafAlloc),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PhysicalLeafPagesInUseBytes,
			prometheus.GaugeValue,
			float64(s.LeafInuse),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Buckets,
			prometheus.GaugeValue,
			float64(s.BucketN),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.InlinedBuckets,
			prometheus.GaugeValue,
			float64(s.InlineBucketN),
			c.name,
			bucket,
		)

		ch <- prometheus.MustNewConstMetric(
			c.InlinedBucketsInUseBytes,
			prometheus.GaugeValue,
			float64(s.InlineBucketInuse),
			c.name,
			bucket,
		)

		return nil
	})
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.Buckets, err)
	}
}
