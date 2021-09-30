package prombolt

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	bolt "go.etcd.io/bbolt"
)

func TestBucketStatsCollector(t *testing.T) {
	tests := []struct {
		name    string
		s       []memoryBucketStats
		matches []string
	}{
		{
			name: "single bucket",
			s: []memoryBucketStats{{
				name: "foo",
				s: bolt.BucketStats{
					BranchPageN:       1,
					BranchOverflowN:   2,
					LeafPageN:         3,
					LeafOverflowN:     4,
					KeyN:              5,
					Depth:             6,
					BranchAlloc:       7,
					BranchInuse:       8,
					LeafAlloc:         9,
					LeafInuse:         10,
					BucketN:           11,
					InlineBucketN:     12,
					InlineBucketInuse: 13,
				},
			}},
			matches: []string{
				`bolt_bucket_logical_branch_pages{bucket="foo",database="test.db"} 1`,
				`bolt_bucket_physical_branch_overflow_pages{bucket="foo",database="test.db"} 2`,
				`bolt_bucket_logical_leaf_pages{bucket="foo",database="test.db"} 3`,
				`bolt_bucket_physical_leaf_overflow_pages{bucket="foo",database="test.db"} 4`,
				`bolt_bucket_keys{bucket="foo",database="test.db"} 5`,
				`bolt_bucket_depth{bucket="foo",database="test.db"} 6`,
				`bolt_bucket_physical_branch_pages_allocated_bytes{bucket="foo",database="test.db"} 7`,
				`bolt_bucket_physical_branch_pages_in_use_bytes{bucket="foo",database="test.db"} 8`,
				`bolt_bucket_physical_leaf_pages_allocated_bytes{bucket="foo",database="test.db"} 9`,
				`bolt_bucket_physical_leaf_pages_in_use_bytes{bucket="foo",database="test.db"} 10`,
				`bolt_bucket_buckets{bucket="foo",database="test.db"} 11`,
				`bolt_bucket_inlined_buckets{bucket="foo",database="test.db"} 12`,
				`bolt_bucket_inlined_buckets_in_use_bytes{bucket="foo",database="test.db"} 13`,
			},
		},
		{
			name: "multiple buckets",
			s: []memoryBucketStats{
				{
					name: "foo",
					s: bolt.BucketStats{
						BranchPageN:       1,
						BranchOverflowN:   2,
						LeafPageN:         3,
						LeafOverflowN:     4,
						KeyN:              5,
						Depth:             6,
						BranchAlloc:       7,
						BranchInuse:       8,
						LeafAlloc:         9,
						LeafInuse:         10,
						BucketN:           11,
						InlineBucketN:     12,
						InlineBucketInuse: 13,
					},
				},
				{
					name: "bar",
					s: bolt.BucketStats{
						BranchPageN:       1,
						BranchOverflowN:   2,
						LeafPageN:         3,
						LeafOverflowN:     4,
						KeyN:              5,
						Depth:             6,
						BranchAlloc:       7,
						BranchInuse:       8,
						LeafAlloc:         9,
						LeafInuse:         10,
						BucketN:           11,
						InlineBucketN:     12,
						InlineBucketInuse: 13,
					},
				},
			},
			matches: []string{
				`bolt_bucket_logical_branch_pages{bucket="foo",database="test.db"} 1`,
				`bolt_bucket_physical_branch_overflow_pages{bucket="foo",database="test.db"} 2`,
				`bolt_bucket_logical_leaf_pages{bucket="foo",database="test.db"} 3`,
				`bolt_bucket_physical_leaf_overflow_pages{bucket="foo",database="test.db"} 4`,
				`bolt_bucket_keys{bucket="foo",database="test.db"} 5`,
				`bolt_bucket_depth{bucket="foo",database="test.db"} 6`,
				`bolt_bucket_physical_branch_pages_allocated_bytes{bucket="foo",database="test.db"} 7`,
				`bolt_bucket_physical_branch_pages_in_use_bytes{bucket="foo",database="test.db"} 8`,
				`bolt_bucket_physical_leaf_pages_allocated_bytes{bucket="foo",database="test.db"} 9`,
				`bolt_bucket_physical_leaf_pages_in_use_bytes{bucket="foo",database="test.db"} 10`,
				`bolt_bucket_buckets{bucket="foo",database="test.db"} 11`,
				`bolt_bucket_inlined_buckets{bucket="foo",database="test.db"} 12`,
				`bolt_bucket_inlined_buckets_in_use_bytes{bucket="foo",database="test.db"} 13`,
				`bolt_bucket_logical_branch_pages{bucket="bar",database="test.db"} 1`,
				`bolt_bucket_physical_branch_overflow_pages{bucket="bar",database="test.db"} 2`,
				`bolt_bucket_logical_leaf_pages{bucket="bar",database="test.db"} 3`,
				`bolt_bucket_physical_leaf_overflow_pages{bucket="bar",database="test.db"} 4`,
				`bolt_bucket_keys{bucket="bar",database="test.db"} 5`,
				`bolt_bucket_depth{bucket="bar",database="test.db"} 6`,
				`bolt_bucket_physical_branch_pages_allocated_bytes{bucket="bar",database="test.db"} 7`,
				`bolt_bucket_physical_branch_pages_in_use_bytes{bucket="bar",database="test.db"} 8`,
				`bolt_bucket_physical_leaf_pages_allocated_bytes{bucket="bar",database="test.db"} 9`,
				`bolt_bucket_physical_leaf_pages_in_use_bytes{bucket="bar",database="test.db"} 10`,
				`bolt_bucket_buckets{bucket="bar",database="test.db"} 11`,
				`bolt_bucket_inlined_buckets{bucket="bar",database="test.db"} 12`,
				`bolt_bucket_inlined_buckets_in_use_bytes{bucket="bar",database="test.db"} 13`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testCollector(t, newMemoryBucketStatsCollector(tt.s))

			for _, m := range tt.matches {
				t.Run(m, func(t *testing.T) {
					if !strings.Contains(got, m) {
						t.Fatalf("output did not contain expected metric: %q", m)
					}
				})
			}
		})
	}
}

type memoryBucketStats struct {
	name string
	s    bolt.BucketStats
}

func newMemoryBucketStatsCollector(stats []memoryBucketStats) prometheus.Collector {
	bs := newBucketStatsCollector("test.db", nil)

	bs.forEach = func(fn forEachBucketStatsFunc) error {
		for _, s := range stats {
			if err := fn(s.name, s.s); err != nil {
				return err
			}
		}

		return nil
	}

	return bs
}
