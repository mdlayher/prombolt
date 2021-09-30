package prombolt

import (
	"github.com/prometheus/client_golang/prometheus"
	bolt "go.etcd.io/bbolt"
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

	var (
		labels = []string{"database"}
	)

	return &statsCollector{
		name: name,
		ss:   ss,

		FreelistFreePages: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "freelist_free_pages"),
			"Number of free pages on the freelist.",
			labels,
			nil,
		),

		FreelistPendingPages: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "freelist_pending_pages"),
			"Number of pending pages on the freelist.",
			labels,
			nil,
		),

		FreelistFreePageAllocatedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "freelist_free_page_allocated_bytes"),
			"Number of bytes allocated in free pages on the freelist.",
			labels,
			nil,
		),

		FreelistInUseBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "freelist_in_use_bytes"),
			"Number of bytes in use by the freelist.",
			labels,
			nil,
		),

		ReadTxTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "read_tx_total"),
			"Total number of started read transactions for the database.",
			labels,
			nil,
		),

		OpenReadTx: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, dbSubsystem, "open_read_tx"),
			"Number of currently open read-only transactions for the database.",
			labels,
			nil,
		),

		TxPagesAllocatedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "pages_allocated_total"),
			"Total number of transaction page allocations.",
			labels,
			nil,
		),

		TxPagesAllocatedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "pages_allocated_bytes_total"),
			"Total number of bytes allocated for transaction pages.",
			labels,
			nil,
		),

		TxCursorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "cursors_total"),
			"Total number of cursors created by transactions",
			labels,
			nil,
		),

		TxNodesAllocatedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_allocated_total"),
			"Total number of nodes allocated by transactions.",
			labels,
			nil,
		),

		TxNodesDereferencedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_dereferenced_total"),
			"Total number of nodes dereferenced by transactions.",
			labels,
			nil,
		),

		TxNodeRebalancesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "node_rebalances_total"),
			"Total number of node rebalances by transactions.",
			labels,
			nil,
		),

		TxNodeRebalanceSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "node_rebalance_seconds_total"),
			"Total amount of time in seconds spent rebalancing nodes by transactions",
			labels,
			nil,
		),

		TxNodesSplitTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_split_total"),
			"Total number of nodes split by transactions.",
			labels,
			nil,
		),

		TxNodesSpilledTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_spilled_total"),
			"Total number of nodes spilled by transactions.",
			labels,
			nil,
		),

		TxNodesSpilledSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "nodes_spilled_seconds_total"),
			"Total amount of time in seconds spent spilling nodes by transactions.",
			labels,
			nil,
		),

		TxWritesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "writes_total"),
			"Total number of writes to disk performed by transactions.",
			labels,
			nil,
		),

		TxWriteSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, txSubsystem, "write_seconds_total"),
			"Total amount of time in seconds spent writing to disk by transactions.",
			labels,
			nil,
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

	ch <- prometheus.MustNewConstMetric(
		c.FreelistFreePages,
		prometheus.GaugeValue,
		float64(s.FreePageN),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.FreelistPendingPages,
		prometheus.GaugeValue,
		float64(s.PendingPageN),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.FreelistFreePageAllocatedBytes,
		prometheus.GaugeValue,
		float64(s.FreeAlloc),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.FreelistInUseBytes,
		prometheus.GaugeValue,
		float64(s.FreelistInuse),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.ReadTxTotal,
		prometheus.CounterValue,
		float64(s.TxN),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.OpenReadTx,
		prometheus.GaugeValue,
		float64(s.OpenTxN),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxPagesAllocatedTotal,
		prometheus.CounterValue,
		float64(s.TxStats.PageCount),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxPagesAllocatedBytesTotal,
		prometheus.CounterValue,
		float64(s.TxStats.PageAlloc),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxCursorsTotal,
		prometheus.CounterValue,
		float64(s.TxStats.CursorCount),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxNodesAllocatedTotal,
		prometheus.CounterValue,
		float64(s.TxStats.NodeCount),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxNodesDereferencedTotal,
		prometheus.CounterValue,
		float64(s.TxStats.NodeDeref),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxNodeRebalancesTotal,
		prometheus.CounterValue,
		float64(s.TxStats.Rebalance),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxNodeRebalanceSecondsTotal,
		prometheus.CounterValue,
		s.TxStats.RebalanceTime.Seconds(),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxNodesSplitTotal,
		prometheus.CounterValue,
		float64(s.TxStats.Split),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxNodesSpilledTotal,
		prometheus.CounterValue,
		float64(s.TxStats.Spill),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxNodesSpilledSecondsTotal,
		prometheus.CounterValue,
		s.TxStats.SpillTime.Seconds(),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxWritesTotal,
		prometheus.CounterValue,
		float64(s.TxStats.Write),
		c.name,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TxWriteSecondsTotal,
		prometheus.CounterValue,
		s.TxStats.WriteTime.Seconds(),
		c.name,
	)
}
