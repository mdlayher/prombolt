prombolt [![Build Status](https://travis-ci.org/mdlayher/prombolt.svg?branch=master)](https://travis-ci.org/mdlayher/prombolt) [![GoDoc](http://godoc.org/github.com/mdlayher/prombolt?status.svg)](http://godoc.org/github.com/mdlayher/prombolt) [![Report Card](https://goreportcard.com/badge/github.com/mdlayher/prombolt)](https://goreportcard.com/report/github.com/mdlayher/prombolt)
====

Package `prombolt` provides a [Prometheus](https://prometheus.io/) metrics
collector for [Bolt](https://github.com/boltdb/bolt) databases.
MIT Licensed.

Usage
-----

Instrumenting your application's Bolt database using `prombolt` is trivial.
Simply wrap the database handle using `prombolt.New` and register it with
Prometheus.

```go
const name = "prombolt.db"

db, err := bolt.Open(name, 0666, nil)
if err != nil {
	log.Fatal(err)
}

// Register prombolt handler with Prometheus
prometheus.MustRegister(prombolt.New(name, db))

mux := http.NewServeMux()
mux.Handle("/", newHandler(db))
// Attach Prometheus metrics handler
mux.Handle("/metrics", prometheus.Handler())

http.ListenAndServe(":8080", mux)
```

At this point, Bolt metrics should be available for Prometheus to scrape from
the `/metrics` endpoint of your service.

```
$ curl -s http://localhost:8080/metrics | grep "bolt" | head -n 9
# HELP bolt_bucket_buckets Number of buckets within a bucket, including the top bucket.
# TYPE bolt_bucket_buckets gauge
bolt_bucket_buckets{bucket="foo",database="promboltd.db"} 1
# HELP bolt_bucket_depth Number of levels in B+ tree for a bucket.
# TYPE bolt_bucket_depth gauge
bolt_bucket_depth{bucket="foo",database="promboltd.db"} 1
# HELP bolt_bucket_inlined_buckets Number of inlined buckets for a bucket.
# TYPE bolt_bucket_inlined_buckets gauge
bolt_bucket_inlined_buckets{bucket="foo",database="promboltd.db"} 1
```
