package prombolt

import (
	"github.com/boltdb/bolt"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &statsCollector{}

// A statsCollector is a prometheus.Collector for Bolt database and transaction
// statistics.
type statsCollector struct {
	name string
	ss   statser

	FreelistFreePages              *prometheus.Desc
	FreelistPendingPages           *prometheus.Desc
	FreelistFreePageAllocatedBytes *prometheus.Desc
	FreelistInUseBytes             *prometheus.Desc
	ReadTxTotal                    *prometheus.Desc
	OpenReadTx                     *prometheus.Desc

	TxPagesAllocatedTotal       *prometheus.Desc
	TxPagesAllocatedBytesTotal  *prometheus.Desc
	TxCursorsTotal              *prometheus.Desc
	TxNodesAllocatedTotal       *prometheus.Desc
	TxNodesDereferencedTotal    *prometheus.Desc
	TxNodeRebalancesTotal       *prometheus.Desc
	TxNodeRebalanceSecondsTotal *prometheus.Desc
	TxNodesSplitTotal           *prometheus.Desc
	TxNodesSpilledTotal         *prometheus.Desc
	TxNodesSpilledSecondsTotal  *prometheus.Desc
	TxWritesTotal               *prometheus.Desc
	TxWriteSecondsTotal         *prometheus.Desc
}

var _ statser = &bolt.DB{}

// A statser is a type that can produce a bolt.Stats struct.
type statser interface {
	Stats() bolt.Stats
}

// newStatsCollector creates a new statsCollector with the specified name and
// statser for retrieving statistics.
func newStatsCollector(name string, ss statser) *statsCollector {
	const (
		dbSubsystem = "db"
		txSubsystem = "tx"
	)

	return &statsCollector{
		name: name,
		ss:   ss,

		FreelistFreePages: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "freelist_free_pages"),
			"Number of free pages on the freelist.",
			nil, prometheus.Labels{"database": name},
		),

		FreelistPendingPages: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "freelist_pending_pages"),
			"Number of pending pages on the freelist.",
			nil, prometheus.Labels{"database": name},
		),

		FreelistFreePageAllocatedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "freelist_free_page_allocated_bytes"),
			"Number of bytes allocated in free pages on the freelist.",
			nil, prometheus.Labels{"database": name},
		),

		FreelistInUseBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "freelist_in_use_bytes"),
			"Number of bytes in use by the freelist.",
			nil, prometheus.Labels{"database": name},
		),

		ReadTxTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "read_tx_total"),
			"Total number of started read transactions for the database.",
			nil, prometheus.Labels{"database": name},
		),

		OpenReadTx: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "open_read_tx"),
			"Number of currently open read-only transactions for the database.",
			nil, prometheus.Labels{"database": name},
		),

		TxPagesAllocatedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "pages_allocated_total"),
			"Total number of transaction page allocations.",
			nil, prometheus.Labels{"database": name},
		),

		TxPagesAllocatedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "pages_allocated_bytes_total"),
			"Total number of bytes allocated for transaction pages.",
			nil, prometheus.Labels{"database": name},
		),

		TxCursorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "cursors_total"),
			"Total number of cursors created by transactions",
			nil, prometheus.Labels{"database": name},
		),

		TxNodesAllocatedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_allocated_total"),
			"Total number of nodes allocated by transactions.",
			nil, prometheus.Labels{"database": name},
		),

		TxNodesDereferencedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_dereferenced_total"),
			"Total number of nodes dereferenced by transactions.",
			nil, prometheus.Labels{"database": name},
		),

		TxNodeRebalancesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "node_rebalances_total"),
			"Total number of node rebalances by transactions.",
			nil, prometheus.Labels{"database": name},
		),

		TxNodeRebalanceSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "node_rebalance_seconds_total"),
			"Total amount of time in seconds spent rebalancing nodes by transactions",
			nil, prometheus.Labels{"database": name},
		),

		TxNodesSplitTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_split_total"),
			"Total number of nodes split by transactions.",
			nil, prometheus.Labels{"database": name},
		),

		TxNodesSpilledTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_spilled_total"),
			"Total number of nodes spilled by transactions.",
			nil, prometheus.Labels{"database": name},
		),

		TxNodesSpilledSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_spilled_seconds_total"),
			"Total amount of time in seconds spent spilling nodes by transactions.",
			nil, prometheus.Labels{"database": name},
		),

		TxWritesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "writes_total"),
			"Total number of writes to disk performed by transactions.",
			nil, prometheus.Labels{"database": name},
		),

		TxWriteSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "write_seconds_total"),
			"Total amount of time in seconds spent writing to disk by transactions.",
			nil, prometheus.Labels{"database": name},
		),
	}
}

var _ prometheus.Collector = &statsCollector{}

// Describe implements the prometheus.Collector interface.
func (c *statsCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.FreelistFreePages,
		c.FreelistPendingPages,
		c.FreelistFreePageAllocatedBytes,
		c.FreelistInUseBytes,
		c.ReadTxTotal,
		c.OpenReadTx,

		c.TxPagesAllocatedTotal,
		c.TxPagesAllocatedBytesTotal,
		c.TxCursorsTotal,
		c.TxNodesAllocatedTotal,
		c.TxNodesDereferencedTotal,
		c.TxNodeRebalancesTotal,
		c.TxNodeRebalanceSecondsTotal,
		c.TxNodesSplitTotal,
		c.TxNodesSpilledTotal,
		c.TxNodesSpilledSecondsTotal,
		c.TxWritesTotal,
		c.TxWriteSecondsTotal,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect implements the prometheus.Collector interface.
func (c *statsCollector) Collect(ch chan<- prometheus.Metric) {
	s := c.ss.Stats()

	ch <- prometheus.MustNewConstMetric(c.FreelistFreePages, prometheus.GaugeValue, float64(s.FreePageN))

	ch <- prometheus.MustNewConstMetric(c.FreelistPendingPages, prometheus.GaugeValue, float64(s.PendingPageN))

	ch <- prometheus.MustNewConstMetric(c.FreelistFreePageAllocatedBytes, prometheus.GaugeValue, float64(s.FreeAlloc))

	ch <- prometheus.MustNewConstMetric(c.FreelistInUseBytes, prometheus.GaugeValue, float64(s.FreelistInuse))

	ch <- prometheus.MustNewConstMetric(c.ReadTxTotal, prometheus.CounterValue, float64(s.TxN))

	ch <- prometheus.MustNewConstMetric(c.OpenReadTx, prometheus.GaugeValue, float64(s.OpenTxN))

	ch <- prometheus.MustNewConstMetric(c.TxPagesAllocatedTotal, prometheus.CounterValue, float64(s.TxStats.PageCount))

	ch <- prometheus.MustNewConstMetric(c.TxPagesAllocatedBytesTotal, prometheus.CounterValue, float64(s.TxStats.PageAlloc))

	ch <- prometheus.MustNewConstMetric(c.TxCursorsTotal, prometheus.CounterValue, float64(s.TxStats.CursorCount))

	ch <- prometheus.MustNewConstMetric(c.TxNodesAllocatedTotal, prometheus.CounterValue, float64(s.TxStats.NodeCount))

	ch <- prometheus.MustNewConstMetric(c.TxNodesDereferencedTotal, prometheus.CounterValue, float64(s.TxStats.NodeDeref))

	ch <- prometheus.MustNewConstMetric(c.TxNodeRebalancesTotal, prometheus.CounterValue, float64(s.TxStats.Rebalance))

	ch <- prometheus.MustNewConstMetric(c.TxNodeRebalanceSecondsTotal, prometheus.CounterValue, s.TxStats.RebalanceTime.Seconds())

	ch <- prometheus.MustNewConstMetric(c.TxNodesSplitTotal, prometheus.CounterValue, float64(s.TxStats.Split))

	ch <- prometheus.MustNewConstMetric(c.TxNodesSpilledTotal, prometheus.CounterValue, float64(s.TxStats.Spill))

	ch <- prometheus.MustNewConstMetric(c.TxNodesSpilledSecondsTotal, prometheus.CounterValue, s.TxStats.SpillTime.Seconds())

	ch <- prometheus.MustNewConstMetric(c.TxWritesTotal, prometheus.CounterValue, float64(s.TxStats.Write))

	ch <- prometheus.MustNewConstMetric(c.TxWriteSecondsTotal, prometheus.CounterValue, s.TxStats.WriteTime.Seconds())
}
