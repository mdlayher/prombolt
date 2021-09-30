package prombolt

import (
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	bolt "go.etcd.io/bbolt"
)

func TestStatsCollector(t *testing.T) {
	tests := []struct {
		s       bolt.Stats
		matches []string
	}{
		{
			s: bolt.Stats{
				FreePageN:     1,
				PendingPageN:  2,
				FreeAlloc:     3,
				FreelistInuse: 4,
				TxN:           5,
				OpenTxN:       6,
				TxStats: bolt.TxStats{
					PageCount:     7,
					PageAlloc:     8,
					CursorCount:   9,
					NodeCount:     10,
					NodeDeref:     11,
					Rebalance:     12,
					RebalanceTime: 13 * time.Second,
					Split:         14,
					Spill:         15,
					SpillTime:     16 * time.Second,
					Write:         17,
					WriteTime:     18 * time.Second,
				},
			},
			matches: []string{
				`bolt_db_freelist_free_pages{database="test.db"} 1`,
				`bolt_db_freelist_pending_pages{database="test.db"} 2`,
				`bolt_db_freelist_free_page_allocated_bytes{database="test.db"} 3`,
				`bolt_db_freelist_in_use_bytes{database="test.db"} 4`,
				`bolt_db_read_tx_total{database="test.db"} 5`,
				`bolt_db_open_read_tx{database="test.db"} 6`,
				`bolt_tx_pages_allocated_total{database="test.db"} 7`,
				`bolt_tx_pages_allocated_bytes_total{database="test.db"} 8`,
				`bolt_tx_cursors_total{database="test.db"} 9`,
				`bolt_tx_nodes_allocated_total{database="test.db"} 10`,
				`bolt_tx_nodes_dereferenced_total{database="test.db"} 11`,
				`bolt_tx_node_rebalances_total{database="test.db"} 12`,
				`bolt_tx_node_rebalance_seconds_total{database="test.db"} 13`,
				`bolt_tx_nodes_split_total{database="test.db"} 14`,
				`bolt_tx_nodes_spilled_total{database="test.db"} 15`,
				`bolt_tx_nodes_spilled_seconds_total{database="test.db"} 16`,
				`bolt_tx_writes_total{database="test.db"} 17`,
				`bolt_tx_write_seconds_total{database="test.db"} 18`,
			},
		},
	}

	for _, tt := range tests {
		got := testCollector(t, newMemoryStatsCollector(tt.s))

		for _, m := range tt.matches {
			t.Run(m, func(t *testing.T) {
				if !strings.Contains(got, m) {
					t.Fatalf("output did not contain expected metric: %q", m)
				}
			})
		}
	}
}

func newMemoryStatsCollector(s bolt.Stats) prometheus.Collector {
	return newStatsCollector("test.db", &memoryStatsCollector{
		s: s,
	})
}

var _ statser = &memoryStatsCollector{}

type memoryStatsCollector struct {
	s bolt.Stats
}

func (s *memoryStatsCollector) Stats() bolt.Stats {
	return s.s
}
