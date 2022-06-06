package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	currentCount = 0

	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_http_hit_total",
			Help: "Total number of http hits.",
		},
	)

	//Single

	postSingleConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_post_config_hit_total",
			Help: "Total number of create config hits.",
		},
	)

	postSingleConfigVerHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_post_config_ver_hit_total",
			Help: "Total number of add new config version hits.",
		},
	)

	getSingleConfigVerHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_get_config_ver_hit_total",
			Help: "Total number of get all config versions hits.",
		},
	)

	getSingleConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_get_config_hit_total",
			Help: "Total number of get one config hits.",
		},
	)

	deleteSingleConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_del_config_hit_total",
			Help: "Total number of delete config hits.",
		},
	)

	//Group

	postGroupConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_post_group_hit_total",
			Help: "Total number of post new group hits.",
		},
	)

	postGroupConfigVerHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_post_group_ver_hit_total",
			Help: "Total number of adding new group version hits.",
		},
	)

	getGroupConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_get_group_hit_total",
			Help: "Total number of get group hits.",
		},
	)

	getGroupConfigLabelHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_get_group_config_hit_total",
			Help: "Total number of get group configs by label hits.",
		},
	)

	deleteGroupConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_del_group_hit_total",
			Help: "Total number of delete group hits.",
		},
	)

	metricsList = []prometheus.Collector{
		postSingleConfigHits, getSingleConfigVerHits, postSingleConfigVerHits, getSingleConfigHits,
		deleteSingleConfigHits, postGroupConfigHits, postGroupConfigVerHits, getGroupConfigHits, deleteGroupConfigHits,
		getGroupConfigLabelHits, httpHits,
	}

	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

// POST/singleConfig/
func countPostSingleConfigHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		postSingleConfigHits.Inc()
		f(w, r)
	}
}

// POST/singleConfig/{id}
func countPostSingleConfigVerHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		postSingleConfigVerHits.Inc()
		f(w, r)
	}
}

// GET/singleConfig/{id}
func countGetSingleConfigVerHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getSingleConfigVerHits.Inc()
		f(w, r)
	}
}

// GET/singleConfig{id}/{ver}
func countGetSingleConfigHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getSingleConfigHits.Inc()
		f(w, r)
	}
}

// DELETE/singleConfig/{id}/{ver}
func countDeleteSingleConfigHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		deleteSingleConfigHits.Inc()
		f(w, r)
	}
}

// POST/groupConfig/
func countPostGroupConfigHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		postGroupConfigHits.Inc()
		f(w, r)
	}
}

// POST/groupConfig/{id}
func countPostGroupConfigVerHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		postGroupConfigVerHits.Inc()
		f(w, r)
	}
}

// GET/groupConfig/{id}/{ver}
func countGetGroupConfigHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getGroupConfigHits.Inc()
		f(w, r)
	}
}

// GET/groupConfig/{id}/{ver}/singleConfig/
func countGetGroupConfigsHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getGroupConfigLabelHits.Inc()
		f(w, r)
	}
}

// GET/groupConfig/{id}/{ver}
func countDeleteGroupConfigHits(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		deleteGroupConfigHits.Inc()
		f(w, r)
	}
}
