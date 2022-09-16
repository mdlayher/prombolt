package prombolt

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	bolt "go.etcd.io/bbolt"
)

// testCollector performs a single metrics collection pass against the input
// prometheus.Collector, and returns a string containing metrics output.
func testCollector(t *testing.T, collector prometheus.Collector) string {
	if err := prometheus.Register(collector); err != nil {
		t.Fatalf("failed to register Prometheus collector: %v", err)
	}
	defer prometheus.Unregister(collector)

	promServer := httptest.NewServer(promhttp.Handler())
	defer promServer.Close()

	resp, err := http.Get(promServer.URL)
	if err != nil {
		t.Fatalf("failed to GET data from prometheus: %v", err)
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read server response: %v", err)
	}

	return string(buf)
}

func TestRegisterer(t *testing.T) {
	reg := prometheus.NewRegistry()
	{
		db := new(bolt.DB)
		if err := reg.Register(New("db_A", db)); err != nil {
			t.Fatal(err)
		}
	}
	{
		db := new(bolt.DB)
		if err := reg.Register(New("db_B", db)); err != nil {
			t.Fatal(err)
		}
	}
}
